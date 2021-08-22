package edgerouter

// stringSlice converts an interface slice to a string slice.
func stringSlice(src []interface{}) []string {
	var dst []string
	for _, v := range src {
		dst = append(dst, v.(string))
	}
	return dst
}

// Convert an empty string to a nill pointer so it can be omitted when
// sending to the API
func emptyStringToNil(value string) *string {
	if value == "" {
		return nil
	} else {
		return &value
	}
}

// Pair to nil emptyStringToNil. Read the value back from the API into an empty string for the terraform state
func nilStringToEmpty(value *string) string {
	if value == nil {
		return ""
	} else {
		return *value
	}
}
