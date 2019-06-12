package main

import (
	"github.com/demeesterdev/terraform-provider-transip/transip"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: transip.Provider,
	})
}
