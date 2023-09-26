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

func resourceGERRuleset() *schema.Resource {
	return &schema.Resource{
		Description: "GER Ruleset is a set of rules and configurations in Squadcast. It allows users to define how alerts are routed to services without the need to set up individual webhooks for each alert source.",

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
			"catch_all_action": {
				Description: "The \"Catch-All Action\", when configured, specifies a fall back service. If none of the defined rules for an incoming event evaluate to true, the incoming event is routed to the Catch-All service, ensuring no events are missed.",
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
		return nil, errors.New(fmt.Sprintf("%s is not a valid alert source name. Find all alert sources supported on Squadcast [here](https://www.squadcast.com/integrations).", alertSourceName))
	}

	d.Set("alert_source", alertSourceName)
	d.Set("alert_source_shortname", alertSourceShortName)
	d.Set("alert_source_version", alertSourceVersion)
	d.Set("ger_id", gerID)

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
		return diag.Errorf("%s is not a valid alert source name. Find all alert sources supported on Squadcast [here](https://www.squadcast.com/integrations).", alertSource)
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
	d.Set("alert_source_shortname", req.AlertSourceShortName)
	d.Set("alert_source_version", req.AlertSourceVersion)

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

	gerRuleset, err := client.GetGERRulesetByAlertSource(ctx, d.Get("ger_id").(string), alertSourceData)
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
		Name:    d.Get("alert_source_shortname").(string),
		Version: d.Get("alert_source_version").(string),
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
