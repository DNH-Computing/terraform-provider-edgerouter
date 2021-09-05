package edgerouter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DNH-Computing/terraform-provider-edgerouter/edgerouter/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func edgeosFirewallResource() *schema.Resource {
	return &schema.Resource{
		Read:   edgeosFirewallRead,
		Create: edgeosFirewallCreate,
		Delete: edgeosFirewallDelete,
		Update: edgeosFirewallUpdate,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// TODO copy-paste this for v6
			// "type": {
			// 	Type:     schema.TypeString,
			// 	Required: true,
			// 	ForceNew: true,
			// 	ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
			// 		v := val.(string)
			// 		if v != "ipv4" && v != "ipv6" {
			// 			errs = append(errs, fmt.Errorf("Type must be either 'ipv4' or 'ipv6'. Got %s instead", v))
			// 		}
			// 		return
			// 	},
			// },
			"default_action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"prioity": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
						},
						"log": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"state": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						// TODO Other port types
						"destination_port_group": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func edgeosFirewallRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	firewallName := d.Get("name").(string)
	firewallPolicyGet := model.Input{
		Get: &model.Root{
			Firewall: &model.Firewall{
				Name: map[string]*model.FirewallPolicy{
					firewallName: nil,
				},
			},
		},
	}
	var firewallPolicyOutput model.Output
	err := client.Post(context.Background(), "/api/edge/batch.json", &firewallPolicyGet, &firewallPolicyOutput)
	if err != nil {
		return err
	}
	if err = model.HandleAPIResponse(firewallPolicyOutput); err != nil {
		return err
	}

	var firewallPolicyfromRouter model.Root
	if err = json.Unmarshal(firewallPolicyOutput.Get, &firewallPolicyfromRouter); err != nil {
		return err
	}
	if firewallPolicyfromRouter.Firewall.Name[firewallName] == nil {
		return fmt.Errorf("Firewall named %s was not found", firewallName)
	}

	d.Set("default_action", firewallPolicyfromRouter.Firewall.Name[firewallName].DefaultAction)
	if err := d.Set("rule", edgeosFirewallCovertToTerraform(firewallPolicyfromRouter.Firewall.Name[firewallName].Rule)); err != nil {
		return err
	}
	return nil
}

func edgeosFirewallSetStruct(d *schema.ResourceData) (*model.Root, error) {
	rules, err := edgeosFirewallConvertFromTerraform(d.Get("rule"))
	if err != nil {
		return nil, err
	}

	firewallName := d.Get("name").(string)
	return &model.Root{
		Firewall: &model.Firewall{
			Name: map[string]*model.FirewallPolicy{
				firewallName: {
					DefaultAction: d.Get("default_action").(string),
					Rule:          rules,
				},
			},
		},
	}, nil
}

func edgeosFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	firewallName := d.Get("name").(string)
	var firewallOutput model.Output
	var firewallInput model.Input
	if root, err := edgeosFirewallSetStruct(d); err != nil {
		return err
	} else {
		firewallInput.Set = root
	}

	if err := client.Post(context.Background(), "/api/edge/batch.json", &firewallInput, &firewallOutput); err != nil {
		return err
	}

	d.SetId(firewallName)
	return model.HandleAPIResponse(firewallOutput)
}

func edgeosFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	var firewallOutput model.Output
	firewallDelete := model.Input{
		Delete: edgeosFirewallDeleteStruct(d),
	}

	if err := client.Post(context.Background(), "/api/edge/batch.json", &firewallDelete, &firewallOutput); err != nil {
		return err
	}

	return model.HandleAPIResponse(firewallOutput)
}

func edgeosFirewallDeleteStruct(d *schema.ResourceData) *model.Root {
	firewallName := d.Get("name").(string)
	return &model.Root{
		Firewall: &model.Firewall{
			Name: map[string]*model.FirewallPolicy{
				firewallName: nil,
			},
		},
	}
}

func edgeosFirewallUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	config.Lock.Lock()
	defer config.Lock.Unlock()
	client := config.Client

	var output model.Output
	input := model.Input{
		Delete: edgeosFirewallDeleteStruct(d),
	}
	if root, err := edgeosFirewallSetStruct(d); err != nil {
		return err
	} else {
		input.Set = root
	}

	err := client.Post(context.Background(), "/api/edge/batch.json", &input, &output)
	if err != nil {
		return nil
	}

	return model.HandleAPIResponse(output)
}

// edgeosFirewallCovertToTerraform convers a map[int]*model.FirewallPolicy rule to a terraform data structure.
func edgeosFirewallCovertToTerraform(response map[int]*model.FirewallPolicyRule) (terraformRules []map[string]interface{}) {
	for prioity, rule := range response {
		r := make(map[string]interface{})
		r["prioity"] = prioity
		r["action"] = rule.Action
		r["log"] = rule.Log
		r["protocol"] = nilStringToEmpty(rule.Protocol)
		r["description"] = nilStringToEmpty(rule.Description)
		r["state"] = edgeosFirewallStateFromMap(rule.State)

		if rule.Destination != nil && rule.Destination.Group != nil {
			r["destination_port_group"] = rule.Destination.Group.PortGroup
		}

		terraformRules = append(terraformRules, r)
	}
	return
}

// edgeosFirewallConvertFromTerraform converts
func edgeosFirewallConvertFromTerraform(terraformRule interface{}) (map[int]*model.FirewallPolicyRule, error) {
	rules := make(map[int]*model.FirewallPolicyRule)
	for _, rawRules := range terraformRule.([]interface{}) {
		rule := rawRules.(map[string]interface{})
		prioity := rule["prioity"].(int)
		action := rule["action"].(string)
		log := rule["log"].(bool)
		protocol := rule["protocol"].(string)
		description := rule["description"].(string)
		states := rule["state"].(*schema.Set).List()
		destinationPortGroup := stringSlice(rule["destination_port_group"].(*schema.Set).List())

		if rules[prioity] != nil {
			return nil, fmt.Errorf("Two rules have the same prioity %d", prioity)
		}

		addDestination := false
		destination := model.FirewallPolicyRuleMatch{
			// TODO only create this if we have to
			Group: &model.FirewallPolicyRuleMatchGroup{},
		}
		if destinationPortGroup != nil && len(destinationPortGroup) > 0 {
			destination.Group.PortGroup = destinationPortGroup
			addDestination = true
		}

		rules[prioity] = &model.FirewallPolicyRule{
			Action:      action,
			Log:         model.FirewallBoolean(log),
			Protocol:    emptyStringToNil(protocol),
			Description: emptyStringToNil(description),
			State:       edgeosFirewallStatesFromStrings(stringSlice(states)),
		}

		if addDestination {
			rules[prioity].Destination = &destination
		}
	}
	return rules, nil
}

func edgeosFirewallStatesFromStrings(enabledStates []string) map[string]*model.FirewallBoolean {
	states := make(map[string]*model.FirewallBoolean)
	for _, state := range enabledStates {
		stateForAPI := model.FirewallBoolean(true)
		states[state] = &stateForAPI
	}
	return states
}

func edgeosFirewallStateFromMap(api map[string]*model.FirewallBoolean) (states []string) {
	for state, enabled := range api {
		if enabled != nil && *enabled {
			states = append(states, state)
		}
	}
	return
}
