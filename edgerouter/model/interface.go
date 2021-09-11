package model

// Interface is the main struct used for creating interfaces on the system.
// Many of the interface types are very complex so are split across
// multiple files
type Interface struct {
	Bridge map[string]*BridgeInterface `json:"bridge,omitempty"`
}

// BridgeInterface represents a Bridge type interface on EdgeOS
type BridgeInterface struct {
	Disabled *PresentMarker `json:"disabled,omitempty"`
	Priority int            `json:"priority"`
}
