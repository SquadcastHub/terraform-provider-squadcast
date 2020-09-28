package squadcast

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider represents a resource provider in Terraform
func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{

			"squadcast_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("squadcast_token", nil),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"squadcast_escalation_policy": dataSourceSquadcastEscalationPolicy(),
			// "squadcast_schedule":          dataSourceSquadcastSchedule(),
			// "squadcast_user":              dataSourceSquadcastUser(),
			// "squadcast_squad":              dataSourceSquadcastSquad(),
			// "squadcast_service": dataSourceService(),
		},

		ResourcesMap: map[string]*schema.Resource{
			// "squadcast_escalation_policy":      resourcesquadcastEscalationPolicy(),
			// "squadcast_maintenance_window":     resourcesquadcastMaintenanceWindow(),
			// "squadcast_schedule":               resourceSquadcastSchedule(),
			"squadcast_service": resourceSquadcastService(),
			// "squadcast_service_integration":    resourcesquadcastServiceIntegration(),
			// "squadcast_squad":                   resourcesquadcastSquad(),
			// "squadcast_user":                   resourcesquadcastUser(),
		},
		ConfigureFunc: providerConfigure,
	}

	return p
}

// TODO: Validate api response status code
// func isErrCode(err error, code int) bool {
// if e, ok := err.(*squadcast.Error); ok && e.ErrorResponse.Response.StatusCode == code {
// 	return true
// }

// return false
// 	return true
// }

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	refreshToken := d.Get("squadcast_token").(string)

	if refreshToken == "" {
		return nil, errors.New("Please provide valid refresh token")
	}

	token, err := getAccessToken(refreshToken)
	if err != nil {
		return nil, errors.New("Unable to fetch access token")
	}

	return Config{
		AccessToken: token,
	}, nil

}
