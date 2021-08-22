package edgerouter

import (
	"encoding/json"
	"fmt"
)

type StatusMessage struct {
	Success int             `json:"success,string"`
	Failure int             `json:"failure,string"`
	Error   json.RawMessage `json:"error"`
}

type Output struct {
	Success bool `json:"success"`
}

type MutationOutput struct {
	Output
	Commit *StatusMessage `json:"COMMIT"`
	Set    *StatusMessage `json:"SET,omitempty"`
	Save   *StatusMessage `json:"SAVE"`
	Delete *StatusMessage `json:"DELETE,omitempty"`
}

// stringSlice converts an interface slice to a string slice.
func stringSlice(src []interface{}) []string {
	var dst []string
	for _, v := range src {
		dst = append(dst, v.(string))
	}
	return dst
}

func handleAPIResponse(response MutationOutput) error {
	// Check for sets
	if response.Set != nil && response.Set.Success != 1 {
		return fmt.Errorf("Error in set command: %s", string(response.Set.Error))
	}

	// Check for deletes
	if response.Delete != nil && response.Delete.Success != 1 {
		return fmt.Errorf("Could not execute delete: %s", string(response.Delete.Error))
	}

	if response.Commit != nil && response.Commit.Failure == 1 {
		return fmt.Errorf("Error committing change: %s", string(response.Commit.Error))
	}

	// Check for saves
	if response.Save != nil && response.Save.Success != 1 {
		return fmt.Errorf("Could not save changes: %s", string(response.Save.Error))
	}

	// TODO need to deep-print the objects so pointers are resolved
	if !response.Success {
		return fmt.Errorf("General error occoured. Not further details are avilable")
	}

	return nil
}
