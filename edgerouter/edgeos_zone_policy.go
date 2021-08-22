package edgerouter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PresentMarker bool

func (marker *PresentMarker) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}

func (marker *PresentMarker) isPresent() bool {
	if marker == nil {
		return true
	} else if *marker {
		return false
	} else {
		panic("Is present callled on non-sentinel value. Something is very wrong here.")
	}
}

func sentinelMarker() *PresentMarker {
	marker := PresentMarker(true)
	return &marker
}

func converToMarker(value bool) *PresentMarker {
	if value {
		var marker PresentMarker
		return &marker
	} else {
		return nil
	}
}

type ZonePolicyInput struct {
	Get    *ZonePolicy `json:"GET,omitempty"` // omitempty to not put anything in the JSON if the field is `nil`
	Set    *ZonePolicy `json:"SET,omitempty"`
	Delete *ZonePolicy `json:"DELETE,omitempty"`
}
type ZonePolicyOutput struct {
	Get json.RawMessage `json:"GET"` // This need to be unmarshalled when we read it because the router 'helpfully' resonds and if it's empty the struct isn't correct for unmarshalling. So we deferr the unmarshalling until we want to read it.
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
	LocalZone     *PresentMarker           `json:"local-zone,omitempty"`
}
type ZoneNodeFrom struct {
	Firewall ZoneNodeFromFirewall `json:"firewall"`
}
type ZoneNodeFromFirewall struct {
	Name     *string `json:"name,omitempty"`
	Ipv6Name *string `json:"ipv6-name,omitempty"`
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
	zonePolicyGet := ZonePolicyInput{
		Get: &ZonePolicy{
			ZonePolicy: &ZonePolicyNode{
				Zone: map[string]*ZoneNode{
					zoneName: nil,
				},
			},
		},
	}
	var zonePolicy ZonePolicyOutput
	err := client.Post(context.Background(), "/api/edge/batch.json", &zonePolicyGet, &zonePolicy)
	if err != nil {
		return err
	}
	if !zonePolicy.Success {
		return fmt.Errorf("Could not read zone policy %s", zoneName)
	}

	zonePolicyFromRouter := ZonePolicy{
		ZonePolicy: &ZonePolicyNode{
			Zone: map[string]*ZoneNode{
				zoneName: &ZoneNode{
					LocalZone: sentinelMarker(),
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
	d.Set("local_zone", zonePolicyFromRouter.ZonePolicy.Zone[zoneName].LocalZone.isPresent())
	return nil
}

func edgeosZonePolicyCreateStruct(d *schema.ResourceData) *ZonePolicy {
	zoneName := d.Get("name").(string)
	return &ZonePolicy{
		ZonePolicy: &ZonePolicyNode{
			Zone: map[string]*ZoneNode{
				zoneName: {
					DefaultAction: d.Get("default_action").(string),
					Interface:     stringSlice(d.Get("interfaces").(*schema.Set).List()),
					LocalZone:     converToMarker(d.Get("local_zone").(bool)),
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
	zonePolicySet := ZonePolicyInput{
		Set: edgeosZonePolicyCreateStruct(d),
	}

	var zonePolicySetOutput ZonePolicyOutput

	err := client.Post(context.Background(), "/api/edge/batch.json", &zonePolicySet, &zonePolicySetOutput)
	if err != nil {
		return err
	}

	d.SetId(zoneName)
	return handleAPIResponse(zonePolicySetOutput.MutationOutput)
}

func edgeosZonePolicyDeleteStruct(d *schema.ResourceData) *ZonePolicy {
	zoneName := d.Get("name").(string)
	return &ZonePolicy{
		ZonePolicy: &ZonePolicyNode{
			Zone: map[string]*ZoneNode{
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

	var zonePolicyDelteOutput ZonePolicyOutput
	zonePolicyDelete := ZonePolicyInput{
		Delete: edgeosZonePolicyDeleteStruct(d),
	}

	err := client.Post(context.Background(), "/api/edge/batch.json", &zonePolicyDelete, &zonePolicyDelteOutput)

	if err != nil {
		return err
	}

	return handleAPIResponse(zonePolicyDelteOutput.MutationOutput)
}

func edgeosZonePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	var zonePolicyDelteOutput ZonePolicyOutput
	zonePolicyDelete := ZonePolicyInput{
		Delete: edgeosZonePolicyDeleteStruct(d),
		Set:    edgeosZonePolicyCreateStruct(d),
	}

	err := client.Post(context.Background(), "/api/edge/batch.json", &zonePolicyDelete, &zonePolicyDelteOutput)

	if err != nil {
		return err
	}

	return handleAPIResponse(zonePolicyDelteOutput.MutationOutput)
}
