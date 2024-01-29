package provider

import (
	"context"
	"fmt"

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
				Type:        schema.TypeSet,
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
				Computed:    true,
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
							Description:  "The type of the maintainer. (user, team or squad)",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"user", "squad", "team"}, false),
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
				Description: "List of active alert source names. Find all alert sources supported on Squadcast [here](https://www.squadcast.com/integrations).",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"active_alert_source_endpoints": {
				Description: "Active alert source endpoints.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"alert_source_endpoints": {
				Description: "All available alert source endpoints.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"slack_channel_id": {
				Description: "Slack extension for the service. If set, specifies the ID of the Slack channel associated with the service. If this ID is set, it cannot be removed, but it can be changed to a different slack_channel_id.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"delay_notification_config": {
				Description: "Delay notification config.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Enable delay notification",
						},
						"timezone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Timezone",
						},
						"fixed_timeslot_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Fixed timeslot configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Description: "Start time for the fixed timeslot",
										Type:        schema.TypeString,
										Required:    true,
									},
									"end_time": {
										Description: "End time for the fixed timeslot",
										Type:        schema.TypeString,
										Required:    true,
									},
									"repeat_days": {
										Description: "Repeat days for the fixed timeslot",
										Type:        schema.TypeSet,
										Required:    true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}, false),
										},
										MinItems: 1,
										MaxItems: 7,
									},
								},
							},
						},
						"custom_timeslots_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enable custom timeslots",
						},
						"custom_timeslots": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Custom timeslots",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day_of_week": {
										Type:         schema.TypeString,
										Required:     true,
										Description:  "Day of the week",
										ValidateFunc: validation.StringInSlice([]string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}, false),
									},
									"start_time": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Start time for the custom timeslot",
									},
									"end_time": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "End time for the custom timeslot",
									},
								},
							},
						},
						"assigned_to": {
							Description: "Assignee details",
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Description:  "The id of the assignee.",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: tf.ValidateObjectID,
									},
									"type": {
										Description:  "The type of the assignee. (user, squad, escalation_policy, service_owner, default_escalation_policy)",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"user", "squad", "escalation_policy", "service_owner", "default_escalation_policy"}, false),
									},
								},
							},
						},
					},
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

	delayNotifConfig := d.Get("delay_notification_config").([]interface{})
	if len(delayNotifConfig) > 0 {
		cfg, err := decodeDelayNotificationConfig(delayNotifConfig)
		if err != nil {
			return err
		}
		serviceCreateReq.DelayNotificationConfig = cfg
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
		for _, malertsource := range malertsources {
			for _, alertSource := range alertSources {
				if alertSource.Type == malertsource {
					alertSourceIDs = append(alertSourceIDs, alertSource.ID)
					break
				}
				if alertSource.Type != malertsource && alertSource.Type == alertSources[len(alertSources)-1].Type {
					return diag.Errorf("%s is not a valid alert source name. Find all alert sources supported on Squadcast on https://www.squadcast.com/integrations", malertsource)
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

	mdependencies := tf.ExpandStringSet(d.Get("dependencies").(*schema.Set))
	if len(mdependencies) > 0 {
		_, err = client.UpdateServiceDependencies(ctx, service.ID, &api.UpdateServiceDependenciesReq{
			Data: mdependencies,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	slackChannelID, exists := d.GetOk("slack_channel_id")
	if exists {
		_, err = client.UpdateSlackChannel(ctx, service.ID, &api.AddSlackChannelReq{
			ChannelID: slackChannelID.(string),
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

	var activeAlertSourcesMap = make(map[string]string, len(activeAlertSources.AlertSources))
	for _, alertSource := range activeAlertSources.AlertSources {
		for _, malertsource := range alertSources {
			if alertSource.ID == malertsource.ID {
				activeAlertSourcesMap[malertsource.ShortName] = malertsource.Endpoint(client.IngestionBaseURL, service)
			}
		}
	}

	service.ActiveAlertSources = activeAlertSourcesMap

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

	delayNotifConfig := d.Get("delay_notification_config").([]interface{})
	if len(delayNotifConfig) > 0 {
		cfg, err := decodeDelayNotificationConfig(delayNotifConfig)
		if err != nil {
			return err
		}
		updateReq.DelayNotificationConfig = cfg
	}

	_, err := client.UpdateService(ctx, d.Id(), &updateReq)
	if err != nil {
		return diag.FromErr(err)
	}

	malertsources := tf.ListToSlice[string](d.Get("alert_sources"))
	if len(malertsources) == 0 && d.HasChange("alert_sources") {
		alertSourcesReq := api.AddAlertSourcesReq{
			AlertSources: []string{},
		}
		_, err = client.AddAlertSources(ctx, d.Id(), &alertSourcesReq)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if len(malertsources) > 0 {
		var alertSourceIDs []string
		alertSources, err := client.ListAlertSources(ctx)
		for _, malertsource := range malertsources {
			for _, alertSource := range alertSources {
				if alertSource.Type == malertsource {
					alertSourceIDs = append(alertSourceIDs, alertSource.ID)
					break
				}
				if alertSource.Type != malertsource && alertSource.Type == alertSources[len(alertSources)-1].Type {
					return diag.Errorf("%s is not a valid alert source name. Find all alert sources supported on Squadcast on https://www.squadcast.com/integrations", malertsource)
				}
			}
		}
		if len(alertSourceIDs) == 0 {
			return diag.Errorf("Invalid alert sources provided.")
		}
		alertSourcesReq := api.AddAlertSourcesReq{
			AlertSources: alertSourceIDs,
		}
		_, err = client.AddAlertSources(ctx, d.Id(), &alertSourcesReq)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	mdependencies := tf.ExpandStringSet(d.Get("dependencies").(*schema.Set))
	if len(mdependencies) > 0 {
		_, err = client.UpdateServiceDependencies(ctx, d.Id(), &api.UpdateServiceDependenciesReq{
			Data: mdependencies,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("slack_channel_id") {
		_, err = client.UpdateSlackChannel(ctx, d.Id(), &api.AddSlackChannelReq{
			ChannelID: d.Get("slack_channel_id").(string),
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

func decodeDelayNotificationConfig(delayNotifConfig []interface{}) (*api.NotificationsDelayConfig, diag.Diagnostics) {
	delayNotifConfigMap, ok := delayNotifConfig[0].(map[string]interface{})
	if !ok {
		return nil, diag.Errorf("delay_notification_config is invalid")
	}

	isCustomTimeSlotEnabled := delayNotifConfigMap["custom_timeslots_enabled"].(bool)
	cfg := &api.NotificationsDelayConfig{
		IsEnabled:              delayNotifConfigMap["is_enabled"].(bool),
		Timezone:               delayNotifConfigMap["timezone"].(string),
		CustomTimeslotsEnabled: isCustomTimeSlotEnabled,
	}

	assignTo := delayNotifConfigMap["assigned_to"].([]interface{})
	if len(assignTo) > 0 {
		assigneeMap, ok := assignTo[0].(map[string]interface{})
		if !ok {
			return nil, diag.Errorf("assigned_to is invalid")
		}

		cfg.AssignedTo = &api.AssignTo{
			ID:   assigneeMap["id"].(string),
			Type: assigneeMap["type"].(string),
		}
	}

	fixedTimeSlot := delayNotifConfigMap["fixed_timeslot_config"].([]interface{})
	if len(fixedTimeSlot) > 0 {
		if isCustomTimeSlotEnabled {
			return nil, diag.Errorf("fixed_timeslot_config and custom_timeslots cannot be enabled at the same time")
		}
		fixedTimeSlotMap, ok := fixedTimeSlot[0].(map[string]interface{})
		if !ok {
			return nil, diag.Errorf("fixed_timeslot_config is invalid")
		}
		repeatDays := fixedTimeSlotMap["repeat_days"].(*schema.Set).List()
		repeatDaysInt := make([]int, 0, len(repeatDays))
		for _, repeatDay := range repeatDays {
			repeatDayStr := repeatDay.(string)
			day := api.DayOfWeekMap[repeatDayStr]
			repeatDaysInt = append(repeatDaysInt, int(day))
		}

		cfg.FixedTimeslotConfig = &api.FixedTimeslotConfig{
			DelayTimeSlot: api.DelayTimeSlot{
				StartTime: fixedTimeSlotMap["start_time"].(string),
				EndTime:   fixedTimeSlotMap["end_time"].(string),
			},
			RepeatOnDays: repeatDaysInt,
		}
	}

	if isCustomTimeSlotEnabled {
		customTimeSlots := delayNotifConfigMap["custom_timeslots"].(*schema.Set).List()
		cfg.CustomTimeslots = make(map[string][]api.DelayTimeSlot)
		if len(customTimeSlots) > 0 {
			for _, customTimeSlot := range customTimeSlots {
				customTimeSlotMap, ok := customTimeSlot.(map[string]interface{})
				if !ok {
					return nil, diag.Errorf("custom_timeslots is invalid")
				}

				var customTimeSlotConfig = api.DelayTimeSlot{
					StartTime: customTimeSlotMap["start_time"].(string),
					EndTime:   customTimeSlotMap["end_time"].(string),
				}

				dayOfWeek := api.DayOfWeekMap[customTimeSlotMap["day_of_week"].(string)]
				cfg.CustomTimeslots[fmt.Sprintf("%d", int(dayOfWeek))] = append(cfg.CustomTimeslots[fmt.Sprintf("%d", int(dayOfWeek))], customTimeSlotConfig)
			}
		} else {
			return nil, diag.Errorf("custom_timeslots cannot be empty when custom_timeslots_enabled is true")
		}
	}
	return cfg, nil
}
