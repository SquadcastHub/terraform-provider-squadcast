package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func dataSourceService() *schema.Resource {
	return &schema.Resource{
		Description: "[Squadcast Services](https://support.squadcast.com/docs/adding-a-service-1) are the core components of your infrastructure/application for which alerts are generated. Services in Squadcast represent specific systems, applications, components, products, or teams for which an incident is created. To check out some of the best practices on creating Services in Squadcast, refer to the guide [here](https://www.squadcast.com/blog/how-to-configure-services-in-squadcast-best-practices-to-reduce-mttr)." +
			"Use this data source to get information about a specific service.",
		ReadContext: dataSourceServiceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Service id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:  "Name of the Service.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": {
				Description: "Detailed description about the service.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"team_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
			},
			"escalation_policy_id": {
				Description: "Escalation policy id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email_prefix": {
				Description: "Email prefix.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"api_key": {
				Description: "Unique API key of the service",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description: "Email.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"dependencies": {
				Description: "dependencies.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: tf.ValidateObjectID,
				},
			},
			"alert_source_endpoints": {
				Description: "alert sources.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceServiceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	name, ok := d.GetOk("name")
	if !ok {
		return diag.Errorf("invalid service name provided")
	}

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	tflog.Info(ctx, "Reading service by name", tf.M{
		"name": name.(string),
	})
	service, err := client.GetServiceByName(ctx, teamID.(string), name.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	alertSources, err := client.ListAlertSources(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	service.AlertSources = alertSources.Available().EndpointMap(client.IngestionBaseURL, service)

	if err = tf.EncodeAndSet(service, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
