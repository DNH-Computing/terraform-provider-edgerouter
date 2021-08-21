package edgerouter

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func edgeosIpv6Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Type:        schema.TypeBool,
			Description: "Enable acquisition of IPv6 address using stateless autoconfig",
		},
	}
}
