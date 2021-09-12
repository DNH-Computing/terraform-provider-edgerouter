package edgerouter

import (
	"context"
	"encoding/json"

	"github.com/DNH-Computing/terraform-provider-edgerouter/edgerouter/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// import (
// 	"fmt"
// 	"regexp"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// )

// func resourceEdgeosInterfaceEthernet() *schema.Resource { return resourceEdgeosInterfaceEthernetV1() }
// func datasourceEdgeosInterfaceEthernet() *schema.Resource {
// 	return datasourceEdgeosInterfaceEthernetV1()
// }

// func datasourceEdgeosInterfaceEthernetV1() *schema.Resource {
// 	return &schema.Resource{
// 		Read: edgeosInterfaceEthernetRead,
// 		Schema: map[string]*schema.Schema{
// 			"interface": {
// 				Type:        schema.TypeString,
// 				Description: "The name of the interface to use, e.g. eth0",
// 				Required:    true,
// 				ForceNew:    true,
// 			},
// 			"address": {
// 				Type:        schema.TypeSet,
// 				Description: "A set of static addresses to assign to the interface",
// 				Required:    false,
// 				Elem: &schema.Schema{
// 					Type: schema.TypeString,
// 				},
// 			},
// 			"address_ipv6_autoconf": {
// 				Type:        schema.TypeBool,
// 				Description: "Enable acquisition of IPv6 address using stateless autoconfig",
// 			},
// 			"address_ipv6_eui64": {
// 				Type:        schema.TypeSet,
// 				Description: "Assign IPv6 address using EUI-64 based on MAC address <h:h:h:h/64>",
// 				Elem: &schema.Schema{
// 					Type: schema.TypeString,
// 				},
// 			},
// 			"ipv6_router_advert_name_server": {
// 				Type:        schame.TypeSet,
// 				Description: "IPv6 DNS Servers to advertise",
// 				Elem: &schema.Schema{
// 					Type: schema.TypeString,
// 				},
// 			},
// 			"ipv6_router_advert_prefix": {
// 				Type:        schema.TypeString,
// 				Description: "IPv6 prefix to advertise out this interface",
// 			},
// 			"ipv6_router_advert_options": {
// 				Type:        schema.TypeArray,
// 				Description: "Extra custom options to add to the radvd configuration",
// 				Elem: &schema.Schema{
// 					Type: schema.TypeString,
// 				},
// 			},
// 			"mtu": {
// 				Type:        schema.TypeInt,
// 				Default:     1500,
// 				Description: "The MTU to set on the interface, must be between 1280 and 9000. Warning: 9k frames are known to crash some routers.",
// 				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
// 					v := val.(int)
// 					if v < 1280 || v > 9000 {
// 						errs = append(errs, fmt.Errorf("%q must be between 1280 and 9000 inclusive, got: %d", key, v))
// 					}
// 					return
// 				},
// 			},
// 			"speed": {
// 				Type:        schema.TypeInt,
// 				Required:    false,
// 				Computed:    true,
// 				Default:     nil,
// 				Description: "Force the interface to a specific speed (10, 100, 1000, or 10000). Set this to undef for auto",
// 				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
// 					v := val.(int)
// 					if v != 10 && v != 100 && v != 1_000 && v != 10_000 && v != nil {
// 						errs = append(errs, fmt.Errorf("%q must be one of 10, 100, 1000, or 10000 or unset for auto. Got %d", key, v))
// 					}
// 					return
// 				},
// 			},
// 			"duplex": {
// 				Type:        schema.TypeBool,
// 				Required:    false,
// 				Default:     nil,
// 				Computed:    true,
// 				Description: "Should duplex be full (true), half (false), or auto (undef)",
// 			},
// 			"bond_group": {
// 				Type:        schema.TypeString,
// 				Required:    false,
// 				Description: "Shuld this interface be added to a bond group such as bond0",
// 				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
// 					v := val.(string)
// 					re := regexp.MustCompile("bond\\d+")
// 					if !re.Match(v) {
// 						errs = append(errs, fmt.Errorf("Bond group interface must be of the form bond\\d+. Got %s", v))
// 					}
// 					return
// 				},
// 			},
// 		},
// 	}
// }

// func resourceEdgeosInterfaceEthernetV1() *schema.Resource {
// 	data_source := *datasourceEdgeosInterfaceEthernetv1()
// 	data_source.TypeSet
// }

// func edgeosInterfaceEthernetRead(d *schema.ResourceData, meta interface{}) error {
// }

func edgeosInterfaceEthernetResource() *schema.Resource {
	return &schema.Resource{
		Create: edgeosInterfaceEthernetCreate,
		Delete: edgeosInterfaceEthernetDelete,
		Update: edgeosInterfaceEthernetUpdate,
		Read:   edgeosInterfaceEthernetRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"configuration": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func edgeosIntefaceEthernetSetNode(d *schema.ResourceData) *model.Root {
	interfaceName := d.Get("name").(string)
	var intefaceConfiguration interface{}
	json.Unmarshal(d.Get("configuration").(string), &interfaceConfiguration)
	return &model.Root{
		Interface: &model.Interface{
			Ethernet: map[string]interface{}{
				interfaceName: ,
			},
		},
	}
}

func edgeosInterfaceEthernetDeleteNode(d *schema.ResourceData) *model.Root {
	interfaceName := d.Get("name").(string)
	return &model.Root{
		Interface: &model.Interface{
			Ethernet: map[string]interface{}{
				interfaceName: nil,
			},
		},
	}
}

func edgeosInterfaceEthernetCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	interfaceName := d.Get("name").(string)
	var output model.Output
	input := model.Input{
		Set: edgeosIntefaceEthernetSetNode(d),
	}

	if err := client.Post(context.Background(), "/api/edge/batch.json", &input, &output); err != nil {
		return err
	}

	if err := model.HandleAPIResponse(output); err != nil {
		return err
	}

	d.SetId(interfaceName)
	return nil
}

func edgeosInterfaceEthernetUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	var output model.Output
	input := model.Input{
		Delete: edgeosInterfaceEthernetDeleteNode(d),
		Set:    edgeosIntefaceEthernetSetNode(d),
	}

	if err := client.Post(context.Background(), "/api/edge/batch.json", &input, &output); err != nil {
		return err
	}

	if err := model.HandleAPIResponse(output); err != nil {
		return err
	}

	return nil
}

func edgeosInterfaceEthernetRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	interfaceName := d.Get("name").(string)
	var output model.Output
	input := model.Input{
		Get: edgeosInterfaceEthernetDeleteNode(d),
	}

	if err := client.Post(context.Background(), "/api/edge/batch.json", &input, &output); err != nil {
		return err
	}

	if err := model.HandleAPIResponse(output); err != nil {
		return err
	}

	var root model.Root
	if err := json.Unmarshal(output.Get, &root); err != nil {
		return err
	}

	return d.Set("configuration", root.Interface.Ethernet[interfaceName])
}

func edgeosInterfaceEthernetDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	var output model.Output
	input := model.Input{
		Delete: edgeosInterfaceEthernetDeleteNode(d),
	}

	if err := client.Post(context.Background(), "/api/edge/batch.json", &input, &output); err != nil {
		return err
	}

	if err := model.HandleAPIResponse(output); err != nil {
		return err
	}

	return nil
}
