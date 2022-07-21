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
				DefaultFunc: schema.EnvDefaultFunc("SQUADCAST_TOKEN", nil),
			},
			"dc": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DC", "US"),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"squadcast_escalation_policy": dataSourceSquadcastEscalationPolicy(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"squadcast_service": resourceSquadcastService(),
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

	dc := d.Get("dc").(string)
	if err := setSquadcastAPIHost(dc); err != nil {
		return nil, err
	}

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

// setSquadcastAPIHost updates `squadcastAPIHost` based on the `DC`
func setSquadcastAPIHost(dc string) error {
	if val, ok := apiEndpoints[dc]; ok {
		squadcastAPIHost = val
		return nil
	}
	return errors.New("Please provide valid DC")
}
