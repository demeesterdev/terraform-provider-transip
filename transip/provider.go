package transip

import (
	"fmt"

	tip "github.com/demeesterdev/terraform-provider-transip/transip/helpers/transip"

	"github.com/transip/gotransip"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TRANSIP_ACCOUNT_NAME", nil),
			},
			"private_key_path": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("TRANSIP_KEY_PATH", nil),
				ConflictsWith: []string{"private_key"},
			},
			"private_key": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("TRANSIP_KEY", nil),
				ConflictsWith: []string{"private_key_path"},
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"transip_domain": dataSourceDomain(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"transip_domain": resourceDomain(),
		},
		ConfigureFunc: providerConfigure,
	}

	return p
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	accountName := d.Get("account_name").(string)
	privateKeyPath := d.Get("private_key_path").(string)
	privateKey := d.Get("private_key").(string)

	if privateKey == "" && privateKeyPath == "" {
		return nil, fmt.Errorf("Could not retrieve private key for account '%s'. provide key via provider.transip.private_key or provider.transip.private_key_path", accountName)
	}

	clientConfig := gotransip.ClientConfig{
		AccountName: accountName,
	}
	if privateKeyPath != "" {
		clientConfig.PrivateKeyPath = privateKeyPath
	} else {
		clientConfig.PrivateKeyBody = []byte(privateKey)
	}

	c, err := tip.NewRetryClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("Error building TransIP API Client: %s", err)
	}

	return &c, nil
}
