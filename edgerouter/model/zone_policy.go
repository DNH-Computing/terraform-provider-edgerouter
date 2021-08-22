package model

// ZonePolicy The root of Zone Policy nodes.
type ZonePolicy struct {
	Zone map[string]*ZoneNode `json:"zone"`
}

// ZoneNode a single zone policy
type ZoneNode struct {
	// TODO local zone how?
	DefaultAction string                   `json:"default-action,omitempty"`
	From          map[string]*ZoneNodeFrom `json:"from,omitempty"`
	Interface     []string                 `json:"interface,omitempty"`
	LocalZone     *PresentMarker           `json:"local-zone,omitempty"`
}

// ZoneNodeFrom is the firewall setting for a zone pair
type ZoneNodeFrom struct {
	Firewall ZoneNodeFromFirewall `json:"firewall"`
}

// ZoneNodeFromFirewall the firewall policies to associate to this pair
type ZoneNodeFromFirewall struct {
	Name     *string `json:"name,omitempty"`
	Ipv6Name *string `json:"ipv6-name,omitempty"`
}
