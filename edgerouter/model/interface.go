package model

// Interface is the main struct used for creating interfaces on the system.
// Many of the interface types are very complex so are split across
// multiple files
type Interface struct {
	Ethernet map[string]interface{}      `json:"ethernet,omitempty"`
	Bridge   map[string]*BridgeInterface `json:"bridge,omitempty"`
	Tunnel   map[string]*TunnelInterface `json:"tunnel,omitempty"`
}

// BridgeInterface represents a Bridge type interface on EdgeOS
type BridgeInterface struct {
	Disabled *PresentMarker `json:"disabled,omitempty"`
	Priority int            `json:"priority,string"`
}

// TunnelInterface represents an overlay tunnel
type TunnelInterface struct {
	Encapsulation string  `json:"encapsulation"`
	LocalIP       string  `json:"local-ip"`
	RemoteIP      string  `json:"remote-ip"`
	BridgeGroup   *string `json:"bridge-group,omitempty"`
}
