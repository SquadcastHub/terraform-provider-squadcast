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

func resourceGERRulesetRulesOrdering() *schema.Resource {
	return &schema.Resource{
		Description: "The ordering of rules within a Ruleset dictates the sequence in which rules are evaluated for an alert source. These rules are evaluated sequentially, starting from the top.",

		CreateContext: resourceGERRulesetRulesOrderingUpdate,
		ReadContext:   resourceGERRulesetRulesOrderingRead,
		UpdateContext: resourceGERRulesetRulesOrderingUpdate,
		DeleteContext: resourceGERRulesetRulesOrderingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGERRulesetRulesOrderingImport,
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
			"ordering": {
				Description: "GER Ruleset Rule Ordering.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
		},
	}
}

func resourceGERRulesetRulesOrderingImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
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

func resourceGERRulesetRulesOrderingRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Reading GER Ruleset Rule", tf.M{
		"id": d.Id(),
	})

	alertSourceData := api.GERAlertSource{
		Name:    d.Get("alert_source_shortname").(string),
		Version: d.Get("alert_source_version").(string),
	}

	gerRulesetRules, err := client.GetGERRulesetByAlertSource(ctx, d.Get("ger_id").(string), alertSourceData)
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	gerRulesetRulesOrdering := &api.GERReorderRulesetRules{
		ID:       gerRulesetRules.ID,
		GER_ID:   gerRulesetRules.GER_ID,
		Ordering: gerRulesetRules.Ordering,
	}
	if err = tf.EncodeAndSet(gerRulesetRulesOrdering, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGERRulesetRulesOrderingUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	req := &api.GERReorderRulesetRulesReq{}

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
		return diag.Errorf("%s is not a valid alert source name. Find all alert sources supported on Squadcast [here](https://www.squadcast.com/integrations).", alertSource)
	}

	ordering := d.Get("ordering").([]interface{})

	orderingList := make([]uint, len(ordering))
	for i, v := range ordering {
		val, err := strconv.ParseUint(v.(string), 10, 64)
		if err != nil {
			return diag.Errorf("Invalid rule id.")
		}
		orderingList[i] = uint(val)
	}
	req.Ordering = orderingList

	tflog.Info(ctx, "Updating GER Ruleset Rule Ordering", tf.M{
		"req": req,
	})

	gerRulesetRulesOrdering, err := client.UpdateGERRulesetRulesOrdering(ctx, d.Get("ger_id").(string), api.GERAlertSource{
		Name:    alertSourceShortName,
		Version: alertSourceVersion,
	}, req)
	if err != nil {
		return diag.FromErr(err)
	}

	id := strconv.FormatUint(uint64(gerRulesetRulesOrdering.ID), 10)
	d.SetId(id)
	d.Set("alert_source_shortname", alertSourceShortName)
	d.Set("alert_source_version", alertSourceVersion)

	return resourceGERRulesetRulesOrderingRead(ctx, d, meta)
}

// set state to null
func resourceGERRulesetRulesOrderingDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
