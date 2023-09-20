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

func resourceGERRuleset() *schema.Resource {
	return &schema.Resource{
		Description: "GER Ruleset resource.",

		CreateContext: resourceGERRulesetCreate,
		ReadContext:   resourceGERRulesetRead,
		UpdateContext: resourceGERRulesetUpdate,
		DeleteContext: resourceGERRulesetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGERRulesetImport,
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
			"catch_all_action": {
				Description: "GER Ruleset catch all action.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceGERRulesetImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*api.Client)
	gerID, alertSourceName, alertSourceVersion, err := parse3PartImportID(d.Id())
	if err != nil {
		return nil, err
	}

	alertSourceData := api.GERAlertSource{
		Name:    alertSourceName,
		Version: alertSourceVersion,
	}
	gerRuleset, err := client.GetGERRulesetById(ctx, gerID, alertSourceData)
	if err != nil {
		return nil, err
	}
	gerRulesetID := strconv.FormatUint(uint64(gerRuleset.ID), 10)
	d.SetId(gerRulesetID)

	getAlertSource := map[string]interface{}{
		"name":    gerRuleset.AlertSourceName,
		"version": gerRuleset.AlertSourceVersion,
	}
	d.Set("alert_source", getAlertSource)

	return []*schema.ResourceData{d}, nil
}

func resourceGERRulesetCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	req := &api.GER_Ruleset{}

	alertSource := d.Get("alert_source").(string)
	alertSources, err := client.ListAlertSources(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	isValidAlertSource := false
	for _, alertSourceData := range alertSources {
		if alertSourceData.Type == alertSource {
			req.AlertSourceShortName = alertSourceData.ShortName
			req.AlertSourceVersion = alertSourceData.Version
			isValidAlertSource = true
			break
		}
	}
	if !isValidAlertSource {
		return diag.Errorf("%s is not a valid alert source name. Navigate to Services -> Select any service -> Click Add Alert Source -> Copy the Alert Source name.", alertSource)
	}

	mcatchAllAction := d.Get("catch_all_action").(map[string]interface{})
	catchAllAction := make(map[string]string, len(*&mcatchAllAction))
	for k, v := range *&mcatchAllAction {
		if k != "route_to" {
			return diag.Errorf("%s is not a valid catch all action. Valid catch_all_actions are: route_to", k)
		}
		catchAllAction[k] = v.(string)
	}
	req.CatchAllAction = catchAllAction

	tflog.Info(ctx, "Creating GER Ruleset", tf.M{
		"req": req,
	})
	gerRuleset, err := client.CreateGERRuleset(ctx, d.Get("ger_id").(string), req)
	if err != nil {
		return diag.FromErr(err)
	}

	gerRulesetID := strconv.FormatUint(uint64(gerRuleset.ID), 10)
	d.SetId(gerRulesetID)

	if err = tf.EncodeAndSet(gerRuleset, d); err != nil {
		return diag.FromErr(err)
	}

	return resourceGERRulesetRead(ctx, d, meta)
}

func resourceGERRulesetRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Reading GER Ruleset", tf.M{
		"id": d.Id(),
	})

	alertSourceData := api.GERAlertSource{
		Name:    d.Get("alert_source_shortname").(string),
		Version: d.Get("alert_source_version").(string),
	}

	gerRuleset, err := client.GetGERRulesetById(ctx, d.Get("ger_id").(string), alertSourceData)
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(gerRuleset, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGERRulesetUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	req := &api.GER_Ruleset{}

	if d.HasChange("alert_source") {
		prevAlertSource, _ := d.GetChange("alert_source")
		prevVal := prevAlertSource.(string)
		d.Set("alert_source", prevVal)

		return diag.Errorf("alert_source can only be set during creation.")
	}

	mcatchAllAction := d.Get("catch_all_action").(map[string]interface{})
	catchAllAction := make(map[string]string, len(*&mcatchAllAction))
	for k, v := range *&mcatchAllAction {
		if k != "route_to" {
			return diag.Errorf("%s is not a valid catch all action. Valid catch_all_actions are: route_to", k)
		}
		catchAllAction[k] = v.(string)
	}
	req.CatchAllAction = catchAllAction

	tflog.Info(ctx, "Updating GER Ruleset", tf.M{
		"id": d.Id(),
	})

	_, err := client.UpdateGERRuleset(ctx, d.Get("ger_id").(string), api.GERAlertSource{
		Name:    req.AlertSourceName,
		Version: req.AlertSourceVersion,
	}, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGERRulesetRead(ctx, d, meta)
}

func resourceGERRulesetDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	alertSource := api.GERAlertSource{
		Name:    d.Get("alert_source_shortname").(string),
		Version: d.Get("alert_source_version").(string),
	}

	_, err := client.DeleteGERRuleset(ctx, d.Get("ger_id").(string), alertSource)
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
