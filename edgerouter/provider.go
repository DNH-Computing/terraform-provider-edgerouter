package edgerouter

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// Provider registration function
func Provider() *schema.Provider {
	return &schema.Provider{}
}
