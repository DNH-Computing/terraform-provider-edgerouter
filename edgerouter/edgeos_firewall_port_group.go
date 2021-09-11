package edgerouter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DNH-Computing/terraform-provider-edgerouter/edgerouter/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func edgeosFirewallPortGroupResource() *schema.Resource {
	return &schema.Resource{
		Read:   edgeosFirewallPortGroupRead,
		Create: edgeosFirewallPortGroupCreate,
		Delete: edgeosFirewallPortGroupDelete,
		Update: edgeosFirewallPortGroupUpdate,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ports": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func edgeosFirewallPortGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	groupName := d.Get("name").(string)
	var groupOutput model.Output
	groupGet := model.Input{
		Get: edgeosFirewallPortGroupDeleteGetStruct(d),
	}

	if err := client.Post(context.Background(), "/api/edge/batch.json", &groupGet, &groupOutput); err != nil {
		return err
	}

	if err := model.HandleAPIResponse(groupOutput); err != nil {
		return err
	}

	var root model.Root
	if err := json.Unmarshal(groupOutput.Get, &root); err != nil {
		return err
	}

	group := root.Firewall.Group.PortGroup[groupName]
	if group == nil {
		return fmt.Errorf("Port Group %s not found", groupName)
	}

	d.Set("description", group.Description)
	d.Set("ports", group.Port)
	return nil
}

func edgeosFirewallPortGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	var output model.Output
	groupDelete := model.Input{
		Delete: edgeosFirewallPortGroupDeleteGetStruct(d),
	}

	if err := client.Post(context.Background(), "/api/edge/batch.json", &groupDelete, &output); err != nil {
		return err
	}

	if err := model.HandleAPIResponse(output); err != nil {
		return err
	}

	return nil
}

func edgeosFirewallPortGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	groupName := d.Get("name").(string)
	var output model.Output
	set := model.Input{
		Set: edgeosFirewallPortGroupSetStruct(d),
	}

	if err := client.Post(context.Background(), "/api/edge/batch.json", &set, &output); err != nil {
		return err
	}

	if err := model.HandleAPIResponse(output); err != nil {
		return err
	}

	d.SetId(groupName)
	return nil
}

func edgeosFirewallPortGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	var output model.Output
	update := model.Input{
		Delete: edgeosFirewallPortGroupDeleteGetStruct(d),
		Set:    edgeosFirewallPortGroupSetStruct(d),
	}

	if err := client.Post(context.Background(), "/api/edge/batch.json", &update, &output); err != nil {
		return err
	}

	if err := model.HandleAPIResponse(output); err != nil {
		return err
	}

	return nil
}

func edgeosFirewallPortGroupSetStruct(d *schema.ResourceData) *model.Root {
	groupName := d.Get("name").(string)
	return &model.Root{
		Firewall: &model.Firewall{
			Group: &model.FirewallCommonElements{
				PortGroup: map[string]*model.FirewallPortGroup{
					groupName: {
						Description: d.Get("description").(string),
						Port:        stringSlice(d.Get("ports").(*schema.Set).List()),
					},
				},
			},
		},
	}
}

func edgeosFirewallPortGroupDeleteGetStruct(d *schema.ResourceData) *model.Root {
	groupName := d.Get("name").(string)
	return &model.Root{
		Firewall: &model.Firewall{
			Group: &model.FirewallCommonElements{
				PortGroup: map[string]*model.FirewallPortGroup{
					groupName: nil,
				},
			},
		},
	}
}
