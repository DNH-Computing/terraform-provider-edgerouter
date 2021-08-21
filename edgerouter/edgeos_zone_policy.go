package edgerouter

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type GetZonePolicyInput struct {
	Get *ZonePolicy `json:"GET,omitempty"` // omitempty to not put anything in the JSON if the field is `nil`
}
type GetZonePolicyOutput struct {
	Get *ZonePolicy `json:"GET"`
	Output
}

type SetZonePolicyInput struct {
	Set *ZonePolicy `json:"SET,omitempty"`
}
type SetZonePolicyOutput struct {
	Set *ZonePolicy `json:"SET"`
	MutationOutput
}

type DeleteZonePolicyInput struct {
	Delete *ZonePolicy `json:"DELETE"`
}
type DeleteZonePolicyOutput struct {
	MutationOutput
}
type ZonePolicy struct {
	ZonePolicy *ZonePolicyNode `json:"zone-policy"`
}
type ZonePolicyNode struct {
	Zone map[string]*ZoneNode `json:"zone"`
}
type ZoneNode struct {
	// TODO local zone how?
	DefaultAction string                   `json:"default-action,omitempty"`
	From          map[string]*ZoneNodeFrom `json:"from,omitempty"`
	Interface     []string                 `json:"interface,omitempty"`
}
type ZoneNodeFrom struct {
	Firewall ZoneNodeFromFirewall `json:"firewall"`
}
type ZoneNodeFromFirewall struct {
	Name     string `json:"name"`
	Ipv6Name string `json:"ipv6-name"`
}

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
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"interfaces"},
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
	zonePolicyGet := GetZonePolicyInput{
		Get: &ZonePolicy{
			ZonePolicy: &ZonePolicyNode{
				Zone: map[string]*ZoneNode{
					zoneName: nil,
				},
			},
		},
	}
	var zonePolicy GetZonePolicyOutput
	err := client.Post(context.Background(), "/api/edge/batch.json", &zonePolicyGet, &zonePolicy)
	if err != nil {
		return err
	}
	if !zonePolicy.Success {
		return fmt.Errorf("Could not read zone policy %s", zoneName)
	}

	d.Set("default_action", zonePolicy.Get.ZonePolicy.Zone[zoneName].DefaultAction)
	d.Set("interfaces", zonePolicy.Get.ZonePolicy.Zone[zoneName].Interface)
	return nil
}

func edgeosZonePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	zoneName := d.Get("name").(string)
	zonePolicySet := SetZonePolicyInput{
		Set: &ZonePolicy{
			ZonePolicy: &ZonePolicyNode{
				Zone: map[string]*ZoneNode{
					zoneName: {
						DefaultAction: d.Get("default_action").(string),
						Interface:     stringSlice(d.Get("interfaces").(*schema.Set).List()),
					},
				},
			},
		},
	}

	var zonePolicySetOutput SetZonePolicyOutput

	err := client.Post(context.Background(), "/api/edge/batch.json", &zonePolicySet, &zonePolicySetOutput)
	if err != nil {
		return err
	}

	return handleAPIResponse(zonePolicySetOutput.MutationOutput)
}

func edgeosZonePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	zoneName := d.Get("name").(string)
	var zonePolicyDelteOutput DeleteZonePolicyOutput
	zonePolicyDelete := DeleteZonePolicyInput{
		Delete: &ZonePolicy{
			ZonePolicy: &ZonePolicyNode{
				Zone: map[string]*ZoneNode{
					zoneName: nil,
				},
			},
		},
	}

	err := client.Post(context.Background(), "/api/edge/batch.json", &zonePolicyDelete, &zonePolicyDelteOutput)

	if err != nil {
		return err
	}

	return handleAPIResponse(zonePolicyDelteOutput.MutationOutput)
}

func edgeosZonePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	err := edgeosZonePolicyDelete(d, meta)
	if err != nil {
		return err
	}

	err = edgeosZonePolicyCreate(d, meta)
	if err != nil {
		return err
	}

	return nil
}
