package edgerouter

import "fmt"

type CommitStatus struct {
	Success int    `json:"success,string"`
	Failure int    `json:"failure,string"`
	Error   string `json:"error"`
}

type SaveStatus struct {
	Success int `json:"success,string"`
}

type Output struct {
	Success bool `json:"success"`
}

type MutationOutput struct {
	Output
	Commit *CommitStatus `json:"COMMIT"`
	Save   *SaveStatus   `json:"SAVE"`
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
	// TODO need to deep-print the objects so pointers are resolved
	if !response.Success {
		return fmt.Errorf("General error occoured. Not further details are avilable")
	}

	if response.Commit.Failure == 1 {
		return fmt.Errorf("Error committing change: %s", response.Commit.Error)
	}

	// There's probably a save.Failure here and an error but I've not seen it
	if response.Save.Success != 1 {
		return fmt.Errorf("Could not save changes")
	}

	return nil
}
