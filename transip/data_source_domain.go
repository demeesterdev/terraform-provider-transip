package transip

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/transip/gotransip"
	"github.com/transip/gotransip/domain"
)

func dataSourceDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDomainRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name_servers": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "name servers used for the domain",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceDomainRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*gotransip.SOAPClient)
	domainName := d.Get("name").(string)
	d.SetId(domainName)

	info, err := domain.GetInfo(c, domainName)
	if err != nil {
		return fmt.Errorf("Could not request domain info [%s]: %s", domainName, err)
	}

	if err := d.Set("name", domainName); err != nil {
		return err
	}

	if err := d.Set("name_servers", flattenNameServers(info.Nameservers)); err != nil {
		return err
	}

	return nil
}
