package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func resourceSlo() *schema.Resource {
	return &schema.Resource{
		Description: "`squadcast_slo` manages an SLO.",

		CreateContext: resourceSloCreate,
		ReadContext:   resourceSloRead,
		UpdateContext: resourceSloUpdate,
		DeleteContext: resourceSloDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the SLO.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:  "The name of the SLO.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": {
				Description: "Description of the SLO.",
				Type:        schema.TypeString,
				Default:     "Slo created from terraform provider",
				Optional:    true,
			},
			"target_slo": {
				Description: "The target SLO for the time period.",
				Type:        schema.TypeFloat,
				Required:    true,
			},
			"service_ids": {
				Description: "Service IDs associated with the SLO." +
					"Only incidents from the associated services can be promoted as SLO violating incident",
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"slis": {
				Description: "List of indentified SLIs for the SLO",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"time_interval_type": {
				Description:  "Type of the SLO. Values can either be \"rolling\" or \"fixed\"",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"rolling", "fixed"}, false),
			},
			"duration_in_days": {
				Description: "Tracks SLO for the last x days. Required only when SLO time interval type set to \"rolling\"",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"start_time": {
				Description:  "SLO start time. Required only when SLO time interval type set to \"fixed\"",
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},
			"end_time": {
				Description:  "SLO end time. Required only when SLO time interval type set to \"fixed\"",
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},
			"rules": {
				Description: "SLO monitoring checks has rules for monitoring any SLO violation(Or warning signs)",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The ID of the monitoring rule",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"slo_id": {
							Description: "The ID of the SLO",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": {
							Description: "The name of monitoring check." +
								"\"Supported values are \"breached_error_budget\", \"unhealthy_slo\"," +
								"\"increased_false_positives\", \"remaining_error_budget\"",
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{"breached_error_budget", "unhealthy_slo",
								"increased_false_positives", "remaining_error_budget"}, false),
						},
						"threshold": {
							Description: "Threshold for the monitoring check" +
								"Only supported for rules name \"increased_false_positives\" and \"remaining_error_budget\"",
							Type:     schema.TypeInt,
							Optional: true,
						},
						"is_checked": {
							Description: "Is checked?",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
				Optional: true,
			},
			"notify": {
				Description: "Notification rules for SLO violation" +
					"User can either choose to create an incident or get alerted via email",
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The ID of the notification rule",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"slo_id": {
							Description: "The ID of the SLO.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"user_ids": {
							Description: "List of user ID's who should be alerted via email.",
							Type:        schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"squad_ids": {
							Description: "List of Squad ID's who should be alerted via email.",
							Type:        schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"service_id": {
							Description:  "The ID of the service in which the user want to create an incident",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: tf.ValidateObjectID,
						},
					},
				},
			},
			"team_id": {
				Description:  "The team which SLO resource belongs to",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
				ForceNew:     true,
			},
		},
	}
}

var alertsMap = map[string]string{"is_breached_err_budget": "breached_error_budget",
	"breached_error_budget":               "is_breached_err_budget",
	"is_unhealthy_slo":                    "unhealthy_slo",
	"unhealthy_slo":                       "is_unhealthy_slo",
	"increased_false_positives_threshold": "increased_false_positives",
	"increased_false_positives":           "increased_false_positives_threshold",
	"remaining_err_budget_threshold":      "remaining_error_budget",
	"remaining_error_budget":              "remaining_err_budget_threshold",
}

func resourceSloCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)
	rules := make([]*api.SloMonitoringCheck, 0)
	notify := make([]*api.SloNotify, 0)
	sloActions := make([]*api.SloAction, 0)

	err := Decode(d.Get("rules"), &rules)
	if err != nil {
		return diag.FromErr(err)
	}

	err = Decode(d.Get("notify"), &notify)
	if err != nil {
		return diag.FromErr(err)
	}

	ownerID := d.Get("team_id").(string)

	sloActions = formatRulesAndNotify(rules, notify, 0)

	tflog.Info(ctx, "Creating Slos", map[string]interface{}{
		"name": d.Get("name").(string),
	})

	slo, err := client.CreateSlo(ctx, client.OrganizationID, ownerID, &api.Slo{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		TargetSlo:           d.Get("target_slo").(float64),
		ServiceIDs:          tf.ListToSlice[string](d.Get("service_ids")),
		Slis:                tf.ListToSlice[string](d.Get("slis")),
		TimeIntervalType:    d.Get("time_interval_type").(string),
		DurationInDays:      d.Get("duration_in_days").(int),
		StartTime:           d.Get("start_time").(string),
		EndTime:             d.Get("end_time").(string),
		SloMonitoringChecks: rules,
		SloActions:          sloActions,
		OwnerID:             ownerID,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	idStr := strconv.FormatUint(uint64(slo.ID), 10)
	d.SetId(idStr)
	return resourceSloRead(ctx, d, meta)
}

func resourceSloRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	sloID, ok := d.GetOk("id")
	if !ok {
		return diag.Errorf("invalid slo id")
	}

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id")
	}

	tflog.Info(ctx, "Reading Slos", map[string]interface{}{
		"id":      d.Id(),
		"team_id": d.Get("team_id").(string),
	})

	slo, err := client.GetSlo(ctx, client.OrganizationID, teamID.(string), sloID.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	for _, alert := range slo.SloMonitoringChecks {
		alert.Name = alertsMap[alert.Name]
	}

	if err = tf.EncodeAndSet(slo, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSloUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)
	var rules []*api.SloMonitoringCheck
	sloActions := make([]*api.SloAction, 0)
	notify := make([]*api.SloNotify, 0)

	err := Decode(d.Get("rules"), &rules)
	if err != nil {
		return diag.FromErr(err)
	}

	err = Decode(d.Get("notify"), &notify)
	if err != nil {
		return diag.FromErr(err)
	}

	sloID, _ := strconv.ParseInt(d.Id(), 10, 32)
	ownerID := d.Get("team_id").(string)

	sloActions = formatRulesAndNotify(rules, notify, sloID)

	tflog.Info(ctx, "Updating Slos", map[string]interface{}{
		"name": d.Get("name").(string),
	})

	id := d.Id()

	tflog.Info(ctx, "Updating Slos", map[string]interface{}{
		"name": d.Get("name").(string),
	})

	_, err = client.UpdateSlo(ctx, client.OrganizationID, ownerID, id, &api.Slo{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		TargetSlo:           d.Get("target_slo").(float64),
		ServiceIDs:          tf.ListToSlice[string](d.Get("service_ids")),
		Slis:                tf.ListToSlice[string](d.Get("slis")),
		TimeIntervalType:    d.Get("time_interval_type").(string),
		DurationInDays:      d.Get("duration_in_days").(int),
		StartTime:           d.Get("start_time").(string),
		EndTime:             d.Get("end_time").(string),
		SloMonitoringChecks: rules,
		SloActions:          sloActions,
		OwnerID:             ownerID,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSloRead(ctx, d, meta)
}

func resourceSloDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Deleting Slos", map[string]interface{}{
		"name": d.Get("name").(string),
	})

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id")
	}

	_, err := client.DeleteSlo(ctx, client.OrganizationID, teamID.(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// formatRulesAndNotify transform the payload into the format expected by the API and terraform state
func formatRulesAndNotify(rules []*api.SloMonitoringCheck, notify []*api.SloNotify, sloID int64) []*api.SloAction {
	sloActions := make([]*api.SloAction, 0)
	for _, alert := range rules {
		alert.Name = alertsMap[alert.Name]
		alert.IsChecked = true
		alert.SloID = sloID
	}

	for _, userID := range notify[0].UserIDs {
		user := &api.SloAction{
			Type:   "USER",
			UserID: userID,
			SloID:  sloID,
		}
		sloActions = append(sloActions, user)
	}

	for _, squadID := range notify[0].SquadIDs {
		user := &api.SloAction{
			Type:   "SQUAD",
			UserID: squadID,
			SloID:  sloID,
		}
		sloActions = append(sloActions, user)
	}

	if notify[0].ServiceID != "" {
		service := &api.SloAction{
			Type:   "SERVICE",
			UserID: notify[0].ServiceID,
			SloID:  sloID,
		}
		sloActions = append(sloActions, service)
	}

	return sloActions
}
