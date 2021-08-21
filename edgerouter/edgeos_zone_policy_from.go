package edgerouter

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func edgeosZonePolocyFromResource() *schema.Resource {
	return &schema.Resource{
		Read:   edgeosZonePolicyFromRead,
		Create: edgeosZonePolicyFromCreate,
		Delete: edgeosZonePolicyFromDelete,
		Update: edgeosZonePolicyFromCreate,
		Schema: map[string]*schema.Schema{
			"from_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"to_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipv6_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func edgeosZonePolicyFromCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	fromZone := d.Get("from_zone").(string)
	toZone := d.Get("to_zone").(string)
	d.SetId(fmt.Sprintf("%s-to-%s", fromZone, toZone))

	zonePolicyFromSet := SetZonePolicyInput{
		Set: &ZonePolicy{
			ZonePolicy: &ZonePolicyNode{
				Zone: map[string]*ZoneNode{
					toZone: {
						From: map[string]*ZoneNodeFrom{
							fromZone: {
								Firewall: ZoneNodeFromFirewall{
									Name:     d.Get("policy").(string),
									Ipv6Name: d.Get("ipv6_policy").(string),
								},
							},
						},
					},
				},
			},
		},
	}

	var zonePolicySetOutput SetZonePolicyOutput

	err := client.Post(context.Background(), "/api/edge/batch.json", &zonePolicyFromSet, &zonePolicySetOutput)
	if err != nil {
		return err
	}

	return handleAPIResponse(zonePolicySetOutput.MutationOutput)
}

func edgeosZonePolicyFromRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	fromZone := d.Get("from_zone").(string)
	toZone := d.Get("to_zone").(string)

	var zonePolicyfrom GetZonePolicyOutput
	zonePolicyFromGet := GetZonePolicyInput{
		Get: &ZonePolicy{
			ZonePolicy: &ZonePolicyNode{
				Zone: map[string]*ZoneNode{
					toZone: {
						From: map[string]*ZoneNodeFrom{
							fromZone: nil,
						},
					},
				},
			},
		},
	}

	err := client.Post(context.Background(), "/api/edge/batch.json", &zonePolicyFromGet, &zonePolicyfrom)
	if err != nil {
		return err
	}

	if !zonePolicyfrom.Success {
		return fmt.Errorf("General error occoured. Not further details are avilable")
	}

	from := zonePolicyfrom.Get.ZonePolicy.Zone[fromZone].From[toZone]
	d.SetId(fmt.Sprintf("%s-to-%s", fromZone, toZone))
	d.Set("policy", from.Firewall.Name)
	d.Set("ipv6_policy", from.Firewall.Ipv6Name)
	return nil
}

func edgeosZonePolicyFromDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	fromZone := d.Get("from_zone").(string)
	toZone := d.Get("to_zone").(string)

	var deleteOutput DeleteZonePolicyOutput
	delete := DeleteZonePolicyInput{
		Delete: &ZonePolicy{
			ZonePolicy: &ZonePolicyNode{
				Zone: map[string]*ZoneNode{
					toZone: {
						From: map[string]*ZoneNodeFrom{
							fromZone: nil,
						},
					},
				},
			},
		},
	}

	err := client.Post(context.Background(), "/api/edge/batch.json", &delete, &deleteOutput)
	if err != nil {
		return err
	}

	return handleAPIResponse(deleteOutput.MutationOutput)
}
