package model

import (
	"encoding/json"
	"fmt"
)

// FirewallBoolean Used to represent the enable/disable option used for firewall options
type FirewallBoolean bool

func (logging *FirewallBoolean) MarshalJSON() ([]byte, error) {
	if logging == nil {
		return json.Marshal(nil)
	}

	if *logging {
		return json.Marshal("enable")
	}

	return json.Marshal("disable")
}

func (logging *FirewallBoolean) UnmarshalJSON(valueFromAPI []byte) error {
	var value string
	if err := json.Unmarshal(valueFromAPI, &value); err != nil {
		return err
	}

	if value == "enable" {
		*logging = true
	} else if value == "disable" {
		*logging = false
	} else {
		return fmt.Errorf("Got %s instead of enable or disable", value)
	}
	return nil
}

// Firewall is the base type of the firewall
type Firewall struct {
	// The IPv4 Firewall Policies
	Name map[string]*FirewallPolicy `json:"name,omitempty"`

	// The IPv6 Firewall Policies
	Ipv6Name map[string]*FirewallPolicy `json:"ipv6-name,omitempty"`

	Group *FirewallCommonElements `json:"group,omitempty"`
}

// FirewallCommonElements are common elements that are referenced in firewall policies
type FirewallCommonElements struct {
	AddressGroup map[string]*FirewallAddressGroup `json:"address-group,omitempty"`
	PortGroup    map[string]*FirewallPortGroup    `json:"port-group,omitempty"`
}

// FirewallAddressGroup is an address group defintion
type FirewallAddressGroup struct {
	Description *string  `json:"description,omitempty"`
	Address     []string `json:"address,omitempty"`
}

// FirewallPortGroup represents a group of firewall ports
type FirewallPortGroup struct {
	Description string   `json:"description,omitempty"`
	Port        []string `json:"port,omitempty"`
}

// FirewallPolicy is a firewall policy
type FirewallPolicy struct {
	DefaultAction string                      `json:"default-action,omitempty"`
	Rule          map[int]*FirewallPolicyRule `json:"rule,omitempty"`
}

// FirewallPolicyRule is a single rule in a policy
type FirewallPolicyRule struct {
	Action      string                      `json:"action"`
	Protocol    *string                     `json:"protocol,omitempty"`
	Log         FirewallBoolean             `json:"log"`
	Source      *FirewallPolicyRuleMatch    `json:"source,omitempty"`
	Destination *FirewallPolicyRuleMatch    `json:"destination,omitempty"`
	Description *string                     `json:"description,omitempty"`
	State       map[string]*FirewallBoolean `json:"state,omitempty"`
	// TODO there's a bunch more stuff here that we don't use
}

// FirewallPolicyRuleMatch is the normal conditions for a firewall rule
type FirewallPolicyRuleMatch struct {
	Address []string                      `json:"address,omitempty"`
	Port    []string                      `json:"port,omitempty"`
	Group   *FirewallPolicyRuleMatchGroup `json:"group,omitempty"`
}

// FirewallPolicyRuleMatchGroup reference things in the FirewallCommonElements type.
type FirewallPolicyRuleMatchGroup struct {
	AddressGroup []string `json:"address-group,omitempty"`
	PortGroup    *string  `json:"port-group,omitempty"`
	// TODO Network Group
}
