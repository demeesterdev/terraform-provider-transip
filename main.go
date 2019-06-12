package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/demeesterdev/terraform-provider-transip/transip"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: transip.Provider,
	})
}