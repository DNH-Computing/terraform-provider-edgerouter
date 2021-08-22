package model

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
}

// FirewallAddressGroup is an address group defintion
type FirewallAddressGroup struct {
	Address []string `json:"address,omitempty"`
}

// FirewallPolicy is a firewall policy
type FirewallPolicy struct {
	DefaultAction string                      `json:"default-action,omitempty"`
	Rule          map[int]*FirewallPolicyRule `json:"rule,omitempty"`
}

// FirewallPolicyRule is a single rule in a policy
type FirewallPolicyRule struct {
	Action      string                   `json:"action"`
	Protocol    *string                  `json:"protocol,omitempty"`
	Log         bool                     `json:"log"`
	Source      *FirewallPolicyRuleMatch `json:"source,omitempty"`
	Destination *FirewallPolicyRuleMatch `json:"destination,omitempty"`
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
	// TODO Port Group
	// TODO Network Group
}
