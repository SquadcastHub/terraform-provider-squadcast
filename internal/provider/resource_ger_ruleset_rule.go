package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func resourceGERRulesetRule() *schema.Resource {
	return &schema.Resource{
		Description: "GER Ruleset Rule resource.",

		CreateContext: resourceGERRulesetRuleCreate,
		ReadContext:   resourceGERRulesetRuleRead,
		UpdateContext: resourceGERRulesetRuleUpdate,
		DeleteContext: resourceGERRulesetRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGERRulesetRuleImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "GER Ruleset id.",
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
				Optional:    true,
			},
			"alert_source": {
				Description: "GER Ruleset alert source.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"alert_source_version": {
				Description: "GER Ruleset alert source version.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"alert_source_shortname": {
				Description: "GER Ruleset alert source shortname.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"expression": {
				Description: "GER Ruleset Rule expression.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"action": {
				Description: "GER Ruleset Rule action.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceGERRulesetRuleImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*api.Client)
	gerID, alertSourceName, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}
	alertSourceShortName, alertSourceVersion := "", ""
	alertSources, err := client.ListAlertSources(ctx)
	if err != nil {
		return nil, err
	}
	isValidAlertSource := false
	for _, alertSourceData := range alertSources {
		if alertSourceData.Type == alertSourceName {
			alertSourceShortName = alertSourceData.ShortName
			alertSourceVersion = alertSourceData.Version
			isValidAlertSource = true
			break
		}
	}
	if !isValidAlertSource {
		return nil, errors.New(fmt.Sprintf("%s is not a valid alert source name. Navigate to Services -> Select any service -> Click Add Alert Source -> Copy the Alert Source name.", alertSourceName))
	}

	d.Set("alert_source", alertSourceName)
	d.Set("alert_source_shortname", alertSourceShortName)
	d.Set("alert_source_version", alertSourceVersion)
	d.Set("ger_id", gerID)

	return []*schema.ResourceData{d}, nil
}
func resourceGERRulesetRuleCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	req := &api.GER_Ruleset_Rules{
		Description: d.Get("description").(string),
		Expression:  d.Get("expression").(string),
	}

	alertSource := d.Get("alert_source").(string)
	alertSources, err := client.ListAlertSources(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	alertSourceShortName, alertSourceVersion := "", ""
	isValidAlertSource := false
	for _, alertSourceData := range alertSources {
		if alertSourceData.Type == alertSource {
			alertSourceShortName = alertSourceData.ShortName
			alertSourceVersion = alertSourceData.Version
			isValidAlertSource = true
			break
		}
	}
	if !isValidAlertSource {
		return diag.Errorf("%s is not a valid alert source name. Navigate to Services -> Select any service -> Click Add Alert Source -> Copy the Alert Source name.", alertSource)
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
		Name:    alertSourceShortName,
		Version: alertSourceVersion,
	}, req)
	if err != nil {
		return diag.FromErr(err)
	}

	gerRulesetRulesID := strconv.FormatUint(uint64(gerRulesetRules.ID), 10)
	d.SetId(gerRulesetRulesID)
	d.Set("alert_source_shortname", alertSourceShortName)
	d.Set("alert_source_version", alertSourceVersion)

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
