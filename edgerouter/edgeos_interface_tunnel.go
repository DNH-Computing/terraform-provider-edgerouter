package edgerouter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DNH-Computing/terraform-provider-edgerouter/edgerouter/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func edgeosInterfaceTunnelResource() *schema.Resource {
	return &schema.Resource{
		Create: edgeosInterfaceTunnelCreate,
		Delete: edgeosInterfaceTunnelDelete,
		Read:   edgeosInterfaceTunnelRead,
		Update: edgeosInterfaceTunnelUpdate,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"encapsulation": {
				Type:     schema.TypeString,
				Required: true,
			},
			"local_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remote_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bridge_group": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bridge": {
							Type:     schema.TypeString,
							Required: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"cost": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func edgeosInterfaceTunnelReadDeleteNode(d *schema.ResourceData) *model.Root {
	tunnelName := d.Get("name").(string)
	return &model.Root{
		Interface: &model.Interface{
			Tunnel: map[string]*model.TunnelInterface{
				tunnelName: nil,
			},
		},
	}
}

func edgeosInterfaceTunnelSetNode(d *schema.ResourceData) *model.Root {
	tunnelName := d.Get("name").(string)
	return &model.Root{
		Interface: &model.Interface{
			Tunnel: map[string]*model.TunnelInterface{
				tunnelName: {
					LocalIP:       d.Get("local_ip").(string),
					RemoteIP:      d.Get("remote_ip").(string),
					Encapsulation: d.Get("encapsulation").(string),
					BridgeGroup:   edgeosBridgeGroupConvertFromTerraform(d.Get("bridge_group").([]interface{})),
				},
			},
		},
	}
}

func edgeosInterfaceTunnelRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	tunnelName := d.Get("name").(string)
	var output model.Output
	input := model.Input{
		Get: edgeosInterfaceTunnelReadDeleteNode(d),
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

	if root.Interface.Tunnel[tunnelName] == nil {
		return fmt.Errorf("Could not find tunnel interface %s", tunnelName)
	}

	d.Set("local_ip", root.Interface.Tunnel[tunnelName].LocalIP)
	d.Set("remote_ip", root.Interface.Tunnel[tunnelName].RemoteIP)
	d.Set("encapsulation", root.Interface.Tunnel[tunnelName].Encapsulation)
	if err := d.Set("bridge_group", edgeosBridgeGroupConvertToTerraform(root.Interface.Tunnel[tunnelName].BridgeGroup)); err != nil {
		return err
	}
	return nil
}

func edgeosInterfaceTunnelCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	tunnelName := d.Get("name").(string)
	var output model.Output
	input := model.Input{
		Set: edgeosInterfaceTunnelSetNode(d),
	}

	if err := client.Post(context.Background(), "/api/edge/batch.json", &input, &output); err != nil {
		return err
	}

	if err := model.HandleAPIResponse(output); err != nil {
		return err
	}

	d.SetId(tunnelName)
	return nil
}

func edgeosInterfaceTunnelDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	var output model.Output
	input := model.Input{
		Delete: edgeosInterfaceTunnelReadDeleteNode(d),
	}

	if err := client.Post(context.Background(), "/api/edge/batch.json", &input, &output); err != nil {
		return err
	}

	return model.HandleAPIResponse(output)
}

func edgeosInterfaceTunnelUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	var output model.Output
	input := model.Input{
		Delete: edgeosInterfaceTunnelReadDeleteNode(d),
		Set:    edgeosInterfaceTunnelSetNode(d),
	}

	if err := client.Post(context.Background(), "/api/edge/batch.json", &input, &output); err != nil {
		return err
	}

	if err := model.HandleAPIResponse(output); err != nil {
		return err
	}

	return nil
}

// edgeosBridgeGroupConvertToTerraform converts a possibly-nil model.BridgeGroup to a terraform data structure.
func edgeosBridgeGroupConvertToTerraform(response *model.BridgeGroup) (terraformBridgeGroups []interface{}) {
	if response == nil {
		return nil
	}

	terraformBridgeGroup := make(map[string]interface{})
	terraformBridgeGroup["bridge"] = nilStringToEmpty(response.Bridge)
	if response.Priority != nil {
		terraformBridgeGroup["priority"] = *response.Priority
	}
	if response.Cost != nil {
		terraformBridgeGroup["cost"] = *response.Cost
	}
	return []interface{}{terraformBridgeGroup}
}

// edgeosBridgeGroupConvertFromTerraform converts terraform configuration to a possibly-nil model.BridgeGroup
func edgeosBridgeGroupConvertFromTerraform(terraformBridgeGroups []interface{}) *model.BridgeGroup {
	if len(terraformBridgeGroups) == 0 {
		return nil
	}

	terraformMap := terraformBridgeGroups[0].(map[string]interface{})
	bridgeGroup := &model.BridgeGroup{
		Bridge: emptyStringToNil(terraformMap["bridge"].(string)),
	}
	if priority, ok := terraformMap["priority"]; ok {
		priorityInt := priority.(int)
		bridgeGroup.Priority = &priorityInt
	}
	if cost, ok := terraformMap["cost"]; ok {
		costInt := cost.(int)
		bridgeGroup.Cost = &costInt
	}

	return bridgeGroup
}
