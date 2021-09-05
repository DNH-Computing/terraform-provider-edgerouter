package model

import "encoding/json"

// Input type to the API.
type Input struct {
	Get    *Root `json:"GET,omitempty"`
	Set    *Root `json:"SET,omitempty"`
	Delete *Root `json:"DELETE,omitempty"`
}

// Output is the response from the API
type Output struct {
	// This need to be unmarshalled when we read it because the router 'helpfully' resonds and if it's empty the struct isn't correct for unmarshalling. So we deferr the unmarshalling until we want to read it.
	Get     json.RawMessage `json:"GET"`
	Success bool            `json:"success"`
	Commit  *Status         `json:"COMMIT"`
	Set     *Status         `json:"SET,omitempty"`
	Save    *Status         `json:"SAVE"`
	Delete  *Status         `json:"DELETE,omitempty"`
}

// Root is the root of the configuration tree
type Root struct {
	ZonePolicy *ZonePolicy `json:"zone-policy,omitempty"`
	Firewall   *Firewall   `json:"firewall,omitempty"`
}
