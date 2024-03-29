package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func resourceGERRulesetRule() *schema.Resource {
	return &schema.Resource{
		Description: "GER Ruleset Rules are a set of conditions defined within a Global Event Ruleset. These rules have expressions, whose evaluation will determine the destination service for the incoming events.",

		CreateContext: resourceGERRulesetRuleCreate,
		ReadContext:   resourceGERRulesetRuleRead,
		UpdateContext: resourceGERRulesetRuleUpdate,
		DeleteContext: resourceGERRulesetRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGERRulesetRuleImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "GER Ruleset rule id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ger_id": {
				Description: "GER id.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "GER Ruleset Rule description.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"alert_source": {
				Description: "An alert source refers to the origin of an event (alert), such as a monitoring tool. These alert sources are associated with specific rules in GER, determining where events from each source should be routed. Find all alert sources supported on Squadcast [here](https://www.squadcast.com/integrations).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"alert_source_version": {
				Description: "Version of the linked alert source.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"alert_source_shortname": {
				Description: "Shortname of the linked alert source.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"expression": {
				Description: "An expression is a single condition or a set of conditions that must be met for the rule to take action, such as routing the incoming event to a specific service.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"action": {
				Description: "Rule Action refers to the designated destination service to which an event should be directed towards, whenever a rule expression is true.",
				Type:        schema.TypeMap,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceGERRulesetRuleImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*api.Client)
	gerID, alertSourceName, ruleID, err := parse3PartImportID(d.Id())
	if err != nil {
		return nil, err
	}
	alertSource, err := api.GetAlertSourceDetailsByName(client, ctx, alertSourceName)
	if err != nil {
		return nil, err
	}

	d.Set("alert_source", alertSourceName)
	d.Set("alert_source_shortname", alertSource.ShortName)
	d.Set("alert_source_version", alertSource.Version)
	d.Set("ger_id", gerID)
	d.SetId(ruleID)

	return []*schema.ResourceData{d}, nil
}
func resourceGERRulesetRuleCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	req := &api.GER_Ruleset_Rules{
		Description: d.Get("description").(string),
		Expression:  d.Get("expression").(string),
	}

	alertSourceName := d.Get("alert_source").(string)
	alertSource, err := api.GetAlertSourceDetailsByName(client, ctx, alertSourceName)
	if err != nil {
		return diag.FromErr(err)
	}

	mAction := d.Get("action").(map[string]interface{})
	action := make(map[string]string, len(*&mAction))
	for k, v := range *&mAction {
		if k != "route_to" {
			return diag.Errorf("%s is not a valid action. Valid actions are: route_to", k)
		}
		action[k] = v.(string)
	}
	req.Action = action

	tflog.Info(ctx, "Creating GER Ruleset Rule", tf.M{})
	gerRulesetRules, err := client.CreateGERRulesetRules(ctx, d.Get("ger_id").(string), api.GERAlertSource{
		Name:    alertSource.ShortName,
		Version: alertSource.Version,
	}, req)
	if err != nil {
		return diag.FromErr(err)
	}

	gerRulesetRulesID := strconv.FormatUint(uint64(gerRulesetRules.ID), 10)
	d.SetId(gerRulesetRulesID)
	d.Set("alert_source_shortname", alertSource.ShortName)
	d.Set("alert_source_version", alertSource.Version)

	return resourceGERRulesetRuleRead(ctx, d, meta)
}

func resourceGERRulesetRuleRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Reading GER Ruleset Rule", tf.M{
		"id": d.Id(),
	})

	alertSourceData := api.GERAlertSource{
		Name:    d.Get("alert_source_shortname").(string),
		Version: d.Get("alert_source_version").(string),
	}

	gerRulesetRules, err := client.GetGERRulesetRulesById(ctx, d.Get("ger_id").(string), d.Id(), alertSourceData)
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(gerRulesetRules, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGERRulesetRuleUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	req := &api.GER_Ruleset_Rules{
		Description: d.Get("description").(string),
		Expression:  d.Get("expression").(string),
	}

	if d.HasChange("alert_source") {
		prevAlertSource, _ := d.GetChange("alert_source")
		prevVal := prevAlertSource.(string)
		d.Set("alert_source", prevVal)

		return diag.Errorf("alert_source can only be set during creation.")
	}

	mAction := d.Get("action").(map[string]interface{})
	action := make(map[string]string, len(*&mAction))
	for k, v := range *&mAction {
		if k != "route_to" {
			return diag.Errorf("%s is not a valid action. Valid actions are: route_to", k)
		}
		action[k] = v.(string)
	}
	req.Action = action

	tflog.Info(ctx, "Updating GER Ruleset Rule", tf.M{
		"id": d.Id(),
	})

	_, err := client.UpdateGERRulesetRules(ctx, d.Get("ger_id").(string), d.Id(), api.GERAlertSource{
		Name:    d.Get("alert_source_shortname").(string),
		Version: d.Get("alert_source_version").(string),
	}, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGERRulesetRuleRead(ctx, d, meta)
}

func resourceGERRulesetRuleDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	alertSource := api.GERAlertSource{
		Name:    d.Get("alert_source_shortname").(string),
		Version: d.Get("alert_source_version").(string),
	}

	_, err := client.DeleteGERRulesetRules(ctx, d.Get("ger_id").(string), d.Id(), alertSource)
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
