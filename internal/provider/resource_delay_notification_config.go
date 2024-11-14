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

func resourceDelayedNotificationConfig() *schema.Resource {
	return &schema.Resource{
		Description: "[Delayed Notifications](https://support.squadcast.com/services/delayed-notifications) postpones notifications outside of business hours and provides a summarized report of pending incidents at the start of the next business day.",

		CreateContext: resourceDelayedNotificationConfigCreateOrUpdate,
		ReadContext:   resourceDelayedNotificationConfigRead,
		UpdateContext: resourceDelayedNotificationConfigCreateOrUpdate,
		DeleteContext: resourceDelayedNotificationConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDelayedNotificationConfigImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"service_id": {
				Description:  "Service ID.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
				ForceNew:     true,
			},
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
							Optional:     true,
							Computed:     true,
							ValidateFunc: tf.ValidateObjectID,
						},
						"type": {
							Description:  "The type of the assignee. (user, squad, escalation_policy, default_escalation_policy)",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"user", "squad", "escalation_policy", "default_escalation_policy"}, false),
						},
					},
				},
			},
		},
	}
}

func resourceDelayedNotificationConfigImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	d.Set("service_id", d.Id())

	return []*schema.ResourceData{d}, nil
}

func resourceDelayedNotificationConfigCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)
	serviceID, ok := d.Get("service_id").(string)
	if !ok {
		return diag.Errorf("invalid service_id")
	}

	cfg, err := decodeDelayNotificationConfig(d)
	if err != nil {
		return err
	}

	if cfg.AssignedTo.Type != "default_escalation_policy" && cfg.AssignedTo.ID == "" {
		return diag.Errorf("assigned_to id is required.")
	}

	tflog.Info(ctx, "Updating delayed notification config for service", tf.M{
		"service_id": serviceID,
	})

	_, errx := client.UpdateDelayedNotificationConfig(ctx, serviceID, cfg)
	if errx != nil {
		return diag.FromErr(errx)
	}

	d.SetId(serviceID)

	return resourceDelayedNotificationConfigRead(ctx, d, meta)
}

func resourceDelayedNotificationConfigRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Reading delayed notification config", tf.M{
		"service_id": d.Id(),
	})

	service, err := client.GetServiceById(ctx, "", d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(service.DelayNotificationConfig, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDelayedNotificationConfigDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.UpdateDelayedNotificationConfig(ctx, d.Get("service_id").(string), &api.NotificationsDelayConfig{
		IsEnabled: false,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func decodeDelayNotificationConfig(d *schema.ResourceData) (*api.NotificationsDelayConfig, diag.Diagnostics) {
	delayNotifConfigMap := map[string]interface{}{
		"custom_timeslots_enabled": d.Get("custom_timeslots_enabled"),
		"assigned_to":              d.Get("assigned_to"),
		"is_enabled":               d.Get("is_enabled"),
		"timezone":                 d.Get("timezone"),
		"fixed_timeslot_config":    d.Get("fixed_timeslot_config"),
		"custom_timeslots":         d.Get("custom_timeslots"),
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
