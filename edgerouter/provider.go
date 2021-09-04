package edgerouter

import (
	"crypto/tls"
	"log"
	"sync"

	"github.com/DNH-Computing/terraform-provider-edgerouter/edgerouter/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
			// "edgeos_interface_ethernet": resourceEdgeosInterfaceEthernet(),
			"edgeos_zone_policy":      edgeosZonePolicyResource(),
			"edgeos_zone_policy_from": edgeosZonePolocyFromResource(),
			"edgeos_firewall":         edgeosFirewallResource(),
		},
		ConfigureFunc: configureProvider,
	}
}

type Config struct {
	Lock *sync.Mutex

	Client *client.Client
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	client, err := createClient(d.Get("edge_os_host").(string))

	if err != nil {
		return nil, err
	}

	return &Config{
		Lock:   &sync.Mutex{},
		Client: client,
	}, nil
}

func createClient(host string) (*client.Client, error) {
	user := "ubnt"

	log.Printf("[DEBUG] Creating client for %s@%s", user, host)
	client, err := client.NewClient(&tls.Config{
		InsecureSkipVerify: true, // FIXME
	}, "https://"+host, user, "ubnt") // FIXME user/pass

	return client, err
}
