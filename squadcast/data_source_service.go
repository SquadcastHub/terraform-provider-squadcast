package squadcast

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServiceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				// Sensitive:   true,
				Description: "Service name",
			},
			"description": {
				Type:        schema.TypeString,
				Default:     "Service created via Terraform provider",
				Description: "Service description",
			},
			"escalation_policy_id": {
				Type:        schema.TypeString,
				Description: "Escalation policy id to be associated with the service",
			},
			"email_prefix": {
				Type: schema.TypeString,
				// Optional: true,
				// Default:  true,
				Description: `Email prefix for the service`,
			},
		},
	}
}

// TODO: Check if user provided a valid service name
func dataSourceServiceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
