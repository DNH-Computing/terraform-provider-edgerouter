package edgerouter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DNH-Computing/terraform-provider-edgerouter/edgerouter/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func edgeosZonePolicyResource() *schema.Resource {
	return &schema.Resource{
		Read:   edgeosZonePolicyRead,
		Create: edgeosZonePolicyCreate,
		Delete: edgeosZonePolicyDelete,
		Update: edgeosZonePolicyUpdate,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"interfaces": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"local_zone": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"default_action": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func edgeosZonePolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	zoneName := d.Get("name").(string)
	zonePolicyGet := model.Input{
		Get: &model.Root{
			ZonePolicy: &model.ZonePolicy{
				Zone: map[string]*model.ZoneNode{
					zoneName: nil,
				},
			},
		},
	}
	var zonePolicy model.Output
	err := client.Post(context.Background(), "/api/edge/batch.json", &zonePolicyGet, &zonePolicy)
	if err != nil {
		return err
	}
	if err = model.HandleAPIResponse(zonePolicy); err != nil {
		return err
	}

	zonePolicyFromRouter := model.Root{
		ZonePolicy: &model.ZonePolicy{
			Zone: map[string]*model.ZoneNode{
				zoneName: &model.ZoneNode{
					LocalZone: model.SentinelMarker(),
				},
			},
		},
	}

	if err = json.Unmarshal(zonePolicy.Get, &zonePolicyFromRouter); err != nil {
		return err
	}

	if zonePolicyFromRouter.ZonePolicy.Zone[zoneName] == nil {
		return fmt.Errorf("Zone %s was not found", zoneName)
	}

	d.Set("default_action", zonePolicyFromRouter.ZonePolicy.Zone[zoneName].DefaultAction)
	d.Set("interfaces", zonePolicyFromRouter.ZonePolicy.Zone[zoneName].Interface)
	d.Set("local_zone", zonePolicyFromRouter.ZonePolicy.Zone[zoneName].LocalZone.IsPresent())
	return nil
}

func edgeosZonePolicyCreateStruct(d *schema.ResourceData) *model.Root {
	zoneName := d.Get("name").(string)
	return &model.Root{
		ZonePolicy: &model.ZonePolicy{
			Zone: map[string]*model.ZoneNode{
				zoneName: {
					DefaultAction: d.Get("default_action").(string),
					Interface:     stringSlice(d.Get("interfaces").(*schema.Set).List()),
					LocalZone:     model.ConvertToMarker(d.Get("local_zone").(bool)),
				},
			},
		},
	}
}

func edgeosZonePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	zoneName := d.Get("name").(string)
	zonePolicySet := model.Input{
		Set: edgeosZonePolicyCreateStruct(d),
	}

	var zonePolicySetOutput model.Output

	err := client.Post(context.Background(), "/api/edge/batch.json", &zonePolicySet, &zonePolicySetOutput)
	if err != nil {
		return err
	}

	d.SetId(zoneName)
	return model.HandleAPIResponse(zonePolicySetOutput)
}

func edgeosZonePolicyDeleteStruct(d *schema.ResourceData) *model.Root {
	zoneName := d.Get("name").(string)
	return &model.Root{
		ZonePolicy: &model.ZonePolicy{
			Zone: map[string]*model.ZoneNode{
				zoneName: nil,
			},
		},
	}
}

func edgeosZonePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	var zonePolicyDelteOutput model.Output
	zonePolicyDelete := model.Input{
		Delete: edgeosZonePolicyDeleteStruct(d),
	}

	err := client.Post(context.Background(), "/api/edge/batch.json", &zonePolicyDelete, &zonePolicyDelteOutput)

	if err != nil {
		return err
	}

	return model.HandleAPIResponse(zonePolicyDelteOutput)
}

func edgeosZonePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	var zonePolicyDelteOutput model.Output
	zonePolicyDelete := model.Input{
		Delete: edgeosZonePolicyDeleteStruct(d),
		Set:    edgeosZonePolicyCreateStruct(d),
	}

	err := client.Post(context.Background(), "/api/edge/batch.json", &zonePolicyDelete, &zonePolicyDelteOutput)

	if err != nil {
		return err
	}

	return model.HandleAPIResponse(zonePolicyDelteOutput)
}
