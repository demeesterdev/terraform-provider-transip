package transip

import (
	"fmt"

	"github.com/demeesterdev/terraform-provider-transip/transip/helpers/transip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/transip/gotransip"
	"github.com/transip/gotransip/domain"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainCreate,
		Read:   resourceDomainRead,
		Update: resourceDomainUpdate,
		Delete: resourceDomainDelete,
		Exists: resourceDomainExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name_servers": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "name servers used for the domain",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceDomainCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*gotransip.SOAPClient)
	domainName := d.Get("name").(string)

	domainStatus, err := domain.CheckAvailability(c, domainName)
	if err != nil {
		return fmt.Errorf("Error checking availability for domain [%s]: %+v", domainName, err)
	}

	if domainStatus != domain.StatusFree {
		return transip.CreateDomainUnavailableError(domainName, domainStatus)
	}

	d.Partial(true)

	newDomain := domain.Domain{
		Name: domainName,
	}
	err = domain.Register(c, newDomain)
	if err != nil {
		return fmt.Errorf("Error registrating domain [%s]: %+v", domainName, err)
	}

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error
		_, err = domain.GetInfo(c, domainName)
		if err != nil {
			err = fmt.Errorf("creating domain [%s]: %+v", domainName, err)
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error creating domain [%s]: %+v", domainName, err)
	}

	d.SetId(domainName)

	return resourceDomainUpdate(d, m)
}

func resourceDomainRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*gotransip.SOAPClient)
	domainName := d.Id()

	info, err := domain.GetInfo(c, domainName)
	if err != nil {
		return fmt.Errorf("Could not request domain info [%s]: %s", domainName, err)
	}

	if err := d.Set("name", domainName); err != nil {
		return err
	}

	nameServers := make([]string, 0)
	if info.Nameservers != nil {
		for _, ns := range info.Nameservers {
			nameServers = append(nameServers, ns.Hostname)
		}
	}
	if err := d.Set("name_servers", nameServers); err != nil {
		return err
	}

	return nil
}

func resourceDomainUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	if d.HasChange("name_servers") {
		if err := updateDomainNameServers(d, m); err != nil {
			return err
		}
		d.SetPartial("name_servers")
	}

	d.Partial(false)

	return resourceDomainRead(d, m)
}

func resourceDomainDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDomainExists(d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*gotransip.SOAPClient)
	domainName := d.Id()

	domainStatus, err := domain.CheckAvailability(c, domainName)
	if err != nil {
		return false, fmt.Errorf("Error checking availability for domain [%s]: %+v", domainName, err)
	}

	switch domainStatus {
	case domain.StatusFree:
		return false, nil
	case domain.StatusInYourAccount:
		return true, nil
	}
	return false, transip.CreateDomainUnavailableError(domainName, domainStatus)
}

func updateDomainNameServers(d *schema.ResourceData, m interface{}) error {
	var err error
	c := m.(*gotransip.SOAPClient)
	domainName := d.Id()

	if err := awaitDomainAction(d, m, ""); err != nil {
		return err
	}

	vRaw := d.Get("name_servers").([]interface{})
	err = domain.SetNameservers(c, domainName, expandNameServers(vRaw))
	if err != nil {
		return err
	}

	err = awaitDomainAction(d, m, "changeNameservers")
	return err
}

func awaitDomainAction(d *schema.ResourceData, m interface{}, action string) error {
	c := m.(*gotransip.SOAPClient)
	domainName := d.Id()

	domainAction, err := domain.GetCurrentDomainAction(c, domainName)
	if err != nil {
		return err
	}

	// no current action
	if domainAction.Name == "" {
		if domainAction.HasFailed {
			return fmt.Errorf("Last domain action failed: %s", domainAction.Message)
		} else {
			return nil
		}
	}

	// no action passed need to wait for current action to finish
	if action == "" {
		action = domainAction.Name
	}

	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		var err error
		domainAction, err := domain.GetCurrentDomainAction(c, domainName)
		if err != nil {
			return resource.RetryableError(fmt.Errorf("Domain action %s failed: %s", action, err))
		}

		if domainAction.Name == action {
			if domainAction.HasFailed {
				return resource.NonRetryableError(fmt.Errorf("Domain action %s failed: %s", action, domainAction.Message))
			} else {
				return resource.RetryableError(fmt.Errorf("Domain action %s running", action))
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
