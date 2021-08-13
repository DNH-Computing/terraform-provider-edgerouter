package edgerouter

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/imdario/mergo"
)

func resourceEdgeosInterfaceEthernet() *schema.Resource { return resourceEdgeosInterfaceEthernetV1() }
func datasourceEdgeosInterfaceEthernet() *schema.Resource {
	return datasourceEdgeosInterfaceEthernetV1()
}

func datasourceEdgeosInterfaceEthernetV1() *schema.Resource {
	return &schema.Resource{
		Read: edgeosInterfaceEthernetRead,
		Schema: map[string]*schema.Schema{
			"interface": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceEdgeosInterfaceEthernetV1() *schema.Resource {
	resource := &schema.Resource{}

	mergo.Merge(resource, datasourceEdgeosInterfaceEthernetV1())

	return resource
}

func edgeosInterfaceEthernetRead(d *schema.ResourceData, meta interface{}) error {
}
