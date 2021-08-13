package edgerouter

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// Provider registration function
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"edge_os_host": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"edgeos_interface_ethernet": resourceEdgeosInterfaceEthernet(),
		},
		ConfigureFunc: configureProvider,
	}
}

type Config struct {
	// the hostname to ssh to
	Host string
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	return &Config{
		Host: d.Get("edge_os_host").(string),
	}, nil
}
