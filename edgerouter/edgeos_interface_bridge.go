package edgerouter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DNH-Computing/terraform-provider-edgerouter/edgerouter/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func edgeosInterfaceBridgeResource() *schema.Resource {
	return &schema.Resource{
		Create: edgeosInterfaceBridgeCreate,
		Delete: edgeosInterfaceBridgeDelete,
		Update: edgeosInterfaceBridgeUpdate,
		Read:   edgeosInterfaceBridgeRead,
		Schema: map[string]*schema.Schema{
			"disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  32768,
			},
		},
	}
}

func edgeosInterfaceBridgeReadDeleteNode(d *schema.ResourceData) *model.Root {
	bridgeName := d.Get("name").(string)
	return &model.Root{
		Interface: &model.Interface{
			Bridge: map[string]*model.BridgeInterface{
				bridgeName: nil,
			},
		},
	}
}

func edgeosInterfaceBridgeSetNode(d *schema.ResourceData) *model.Root {
	bridgeName := d.Get("name").(string)
	return &model.Root{
		Interface: &model.Interface{
			Bridge: map[string]*model.BridgeInterface{
				bridgeName: &model.BridgeInterface{
					Disabled: model.ConvertToMarker(d.Get("disabled").(bool)),
					Priority:  d.Get("priority").(int),
				},
			},
		},
	}
}

func edgeosInterfaceBridgeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	var output model.Output
	input := model.Input{
		Get: edgeosInterfaceBridgeReadDeleteNode(d),
	}

	if err := client.Post(context.Background(), "/api/edge/batch.json", &input, &output); err != nil {
		return err
	}

	if err := model.HandleAPIResponse(output); err != nil {
		return err
	}

	bridgeName := d.Get("name").(string)
	root := &model.Root{
		Interface: &model.Interface{
			Bridge: map[string]*model.BridgeInterface{
				bridgeName: {
					Disabled: model.SentinelMarker(),
				},
			},
		},
	}
	if err := json.Unmarshal(output.Get, &root); err != nil {
		return err
	}

	if root.Interface.Bridge[bridgeName] == nil {
		return fmt.Errorf("Could not find bridge interface %v", bridgeName)
	}

	d.Set("disabled", root.Interface.Bridge[bridgeName].Disabled.IsPresent())
	d.Set("priority", root.Interface.Bridge[bridgeName].Priority)
	return nil
}

func edgeosInterfaceBridgeDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	var output model.Output
	input := model.Input{
		Delete: edgeosInterfaceBridgeReadDeleteNode(d),
	}

	if err := client.Post(context.Background(), "/api/edge/batch.json", &input, &output); err != nil {
		return err
	}

	return model.HandleAPIResponse(output)
}

func edgeosInterfaceBridgeCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	var output model.Output
	input := model.Input{
		Set: edgeosInterfaceBridgeSetNode(d),
	}

	if err := client.Post(context.Background(), "/api/edge/batch.json", &input, &output); err != nil {
		return err
	}

	bridgeName := d.Get("name").(string)
	d.SetId(bridgeName)
	return model.HandleAPIResponse(output)
}

func edgeosInterfaceBridgeUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	var output model.Output
	input := model.Input{
		Set: edgeosInterfaceBridgeSetNode(d),
	}

	if err := client.Post(context.Background(), "/api/edge/batch.json", &input, &output); err != nil {
		return err
	}

	bridgeName := d.Get("name").(string)
	d.SetId(bridgeName)
	return model.HandleAPIResponse(output)
}
