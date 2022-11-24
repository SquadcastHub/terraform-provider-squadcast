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

func resourceService() *schema.Resource {
	return &schema.Resource{
		Description: "[Squadcast Services](https://support.squadcast.com/docs/adding-a-service-1) are the core components of your infrastructure/application for which alerts are generated. Services in Squadcast represent specific systems, applications, components, products, or teams for which an incident is created. To check out some of the best practices on creating Services in Squadcast, refer to the guide [here](https://www.squadcast.com/blog/how-to-configure-services-in-squadcast-best-practices-to-reduce-mttr).",

		CreateContext: resourceServiceCreate,
		ReadContext:   resourceServiceRead,
		UpdateContext: resourceServiceUpdate,
		DeleteContext: resourceServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServiceImport,
		},

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
				Description:  "Detailed description about this service.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"team_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
				ForceNew:     true,
			},
			"escalation_policy_id": {
				Description:  "Escalation policy id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
			},
			"email_prefix": {
				Description: "Email prefix.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"api_key": {
				Description: "Unique API key of this service.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description: "Email.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"dependencies": {
				Description: "Dependencies (serviceIds)",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: tf.ValidateObjectID,
				},
			},
			"maintainer": {
				Description: "Service owner.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description:  "The id of the maintainer.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tf.ValidateObjectID,
						},
						"type": {
							Description:  "The type of the maintainer. (user or squad)",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"user", "squad"}, false),
						},
					},
				},
			},
			"tags": {
				Description: "Service tags.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Description: "key",
							Type:        schema.TypeString,
							Required:    true,
						},
						"value": {
							Description: "value",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"alert_sources": {
				Description: "List of alert source names.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"alert_source_endpoints": {
				Description: "Alert source endpoints.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceServiceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	teamID, id, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("team_id", teamID)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func resourceServiceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Creating service", tf.M{
		"name": d.Get("name").(string),
	})
	serviceCreateReq := api.CreateServiceReq{
		Name:               d.Get("name").(string),
		TeamID:             d.Get("team_id").(string),
		Description:        d.Get("description").(string),
		EscalationPolicyID: d.Get("escalation_policy_id").(string),
		EmailPrefix:        d.Get("email_prefix").(string),
	}

	mtags := d.Get("tags").([]any)

	if len(mtags) > 0 {
		var tags []api.ServiceTag
		err := Decode(mtags, &tags)
		if err != nil {
			return diag.FromErr(err)
		}

		serviceCreateReq.Tags = tags
	}

	mmaintainer := d.Get("maintainer").([]interface{})
	if len(mmaintainer) > 0 {
		maintainerMap, ok := mmaintainer[0].(map[string]interface{})
		if !ok {
			return diag.Errorf("maintainer is invalid")
		}

		var maintainer api.ServiceMaintainer
		maintainer.ID = maintainerMap["id"].(string)
		maintainer.Type = maintainerMap["type"].(string)

		serviceCreateReq.Maintainer = &maintainer
	}

	service, err := client.CreateService(ctx, &serviceCreateReq)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(service.ID)

	malertsources := tf.ListToSlice[string](d.Get("alert_sources"))
	if len(malertsources) > 0 {
		var alertSourceIDs []string
		alertSources, err := client.ListAlertSources(ctx)
		for _, alertSource := range alertSources {
			for _, malertsource := range malertsources {
				if alertSource.Type == malertsource {
					alertSourceIDs = append(alertSourceIDs, alertSource.ID)
				}
			}
		}
		if len(alertSourceIDs) == 0 {
			return diag.Errorf("Invalid alert sources provided.")
		}
		alertSourcesReq := api.AddAlertSourcesReq{
			AlertSources: alertSourceIDs,
		}
		_, err = client.AddAlertSources(ctx, service.ID, &alertSourcesReq)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	mdependencies := tf.ListToSlice[string](d.Get("dependencies"))
	if len(mdependencies) > 0 {
		_, err = client.UpdateServiceDependencies(ctx, service.ID, &api.UpdateServiceDependenciesReq{
			Data: mdependencies,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceServiceRead(ctx, d, meta)
}

func resourceServiceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	id := d.Id()

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	tflog.Info(ctx, "Reading service", tf.M{
		"id":   d.Id(),
		"name": d.Get("name").(string),
	})
	service, err := client.GetServiceById(ctx, teamID.(string), id)
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	activeAlertSources, err := client.ListActiveAlertSources(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	alertSources, err := client.ListAlertSources(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var alertSourceNames []string
	for _, alertSource := range activeAlertSources.AlertSources {
		for _, malertsource := range alertSources {
			if alertSource.ID == malertsource.ID {
				alertSourceNames = append(alertSourceNames, malertsource.Type)
			}
		}
	}
	service.ActiveAlertSources = alertSourceNames

	service.AlertSources = alertSources.Available().EndpointMap(client.IngestionBaseURL, service)

	if err = tf.EncodeAndSet(service, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceServiceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	updateReq := api.UpdateServiceReq{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		EscalationPolicyID: d.Get("escalation_policy_id").(string),
		EmailPrefix:        d.Get("email_prefix").(string),
	}

	mtags := d.Get("tags").([]any)

	if len(mtags) > 0 {
		var tags []api.ServiceTag
		err := Decode(mtags, &tags)
		if err != nil {
			return diag.FromErr(err)
		}

		updateReq.Tags = tags
	}

	mmaintainer := d.Get("maintainer").([]interface{})
	if len(mmaintainer) > 0 {
		maintainerMap, ok := mmaintainer[0].(map[string]interface{})
		if !ok {
			return diag.Errorf("maintainer is invalid")
		}

		var maintainer api.ServiceMaintainer
		maintainer.ID = maintainerMap["id"].(string)
		maintainer.Type = maintainerMap["type"].(string)

		updateReq.Maintainer = &maintainer
	}

	_, err := client.UpdateService(ctx, d.Id(), &updateReq)
	if err != nil {
		return diag.FromErr(err)
	}

	malertsources := tf.ListToSlice[string](d.Get("alert_sources"))
	if len(malertsources) > 0 {
		var alertSourceIDs []string
		alertSources, err := client.ListAlertSources(ctx)
		for _, alertSource := range alertSources {
			for _, malertsource := range malertsources {
				if alertSource.Type == malertsource {
					alertSourceIDs = append(alertSourceIDs, alertSource.ID)
				}
			}
		}
		if len(alertSourceIDs) == 0 {
			return diag.Errorf("Invalid alert sources provided")
		}
		alertSourcesReq := api.AddAlertSourcesReq{
			AlertSources: alertSourceIDs,
		}
		_, err = client.AddAlertSources(ctx, d.Id(), &alertSourcesReq)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	mdependencies := tf.ListToSlice[string](d.Get("dependencies"))
	if len(mdependencies) > 0 {
		_, err = client.UpdateServiceDependencies(ctx, d.Id(), &api.UpdateServiceDependenciesReq{
			Data: mdependencies,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceServiceRead(ctx, d, meta)
}

func resourceServiceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteService(ctx, d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
