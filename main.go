package main

import (
	"github.com/DNH-Computing/terraform-provider-edgerouter/edgerouter"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: edgerouter.Provider,
	})
}
