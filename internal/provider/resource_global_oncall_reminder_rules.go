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

func resourceGlobalOncallReminderRules() *schema.Resource {
	return &schema.Resource{
		Description: "[Global Oncall Reminder Rules](https://support.squadcast.com/docs/) implements a global setting for on-call reminder rules to ensure adherence to internal policies and SLAs within a team.",

		CreateContext: resourceGlobalOncallReminderRulesCreate,
		ReadContext:   resourceGlobalOncallReminderRulesRead,
		UpdateContext: resourceGlobalOncallReminderRulesUpdate,
		DeleteContext: resourceGlobalOncallReminderRulesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGlobalOncallReminderRulesImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"team_id": {
				Description:  "Team ID.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
				ForceNew:     true,
			},
			"is_enabled": {
				Description: "Determines whether this setting needs to be enabled or not. When not enabled, each user of the team is expected to set up their own on-call reminder rules.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"rules": {
				Description: "List of on-call reminder rules for the team.",
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description:  "medium of notification. Supported values are 'Email' & 'Push'",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Email", "Push"}, false),
						},
						"time": {
							Description:  "time (in minutes) when the notification needs to be sent before the start of an on-call shift. Max value is 10080 minutes (7 days).",
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtMost(7 * 24 * 60),
						},
					},
				},
			},
		},
	}
}

func resourceGlobalOncallReminderRulesImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	d.Set("team_id", d.Id())

	return []*schema.ResourceData{d}, nil
}

func resourceGlobalOncallReminderRulesCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	req := api.GlobalOncallReminderRulesReq{
		IsEnabled: d.Get("is_enabled").(bool),
		TeamID:    d.Get("team_id").(string),
	}

	rules, errx := decodeReminderRules(d.Get("rules").([]interface{}))
	if errx != nil {
		return errx
	}

	req.Rules = rules

	tflog.Info(ctx, "Creating global oncall reminder rules", tf.M{
		"team_id": d.Get("team_id").(string),
	})

	oncallReminderRules, err := client.CreateGlobalOncallReminderRules(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(oncallReminderRules.TeamID)

	return resourceGlobalOncallReminderRulesRead(ctx, d, meta)
}

func resourceGlobalOncallReminderRulesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	tflog.Info(ctx, "Reading global oncall reminder rules", tf.M{
		"team_id": d.Get("team_id").(string),
	})

	globalOncallReminderRules, err := client.GetGlobalOncallReminderRules(ctx, teamID.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(globalOncallReminderRules, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGlobalOncallReminderRulesUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	req := api.GlobalOncallReminderRulesReq{
		IsEnabled: d.Get("is_enabled").(bool),
	}

	rules, errx := decodeReminderRules(d.Get("rules").([]interface{}))
	if errx != nil {
		return errx
	}

	req.Rules = rules

	_, err := client.UpdateGlobalOncallReminderRules(ctx, d.Get("team_id").(string), &req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGlobalOncallReminderRulesRead(ctx, d, meta)
}

func resourceGlobalOncallReminderRulesDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteGlobalOncallReminderRules(ctx, d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func decodeReminderRules(rules []interface{}) ([]*api.NotificationRule, diag.Diagnostics) {
	rulesReq := []*api.NotificationRule{}

	for _, r := range rules {
		rule, ok := r.(map[string]interface{})
		if !ok {
			return nil, diag.Errorf("invalid rules format")
		}

		rulesReq = append(rulesReq, &api.NotificationRule{
			Time:               rule["time"].(int),
			TypeOfNotification: rule["type"].(string),
		})
	}
	return rulesReq, nil
}
