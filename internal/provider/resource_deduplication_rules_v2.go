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

func resourceDeduplicationRuleV2() *schema.Resource {
	return &schema.Resource{
		Description: "[Deduplication rules](https://support.squadcast.com/docs/de-duplication-rules) can help you reduce alert noise by organising and grouping alerts. This also provides easy access to similar alerts when needed. When these rules evaluate to true for an incoming incident, alerts will get deduplicated.",

		CreateContext: resourceDeduplicationRuleCreateV2,
		ReadContext:   resourceDeduplicationRuleReadV2,
		UpdateContext: resourceDeduplicationRuleUpdateV2,
		DeleteContext: resourceDeduplicationRuleDeleteV2,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDeduplicationRuleImportV2,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"service_id": {
				Description:  "Service id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
				ForceNew:     true,
			},
			"is_basic": {
				Description: "is_basic will be true when users use the drop down selectors which will have lhs, op & rhs value, whereas it will be false when they use the advanced mode and it would have the expression for it's value",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"description": {
				Description: "description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"expression": {
				Description: "The expression which needs to be evaluated to be true for this rule to apply.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"dependency_deduplication": {
				Description: "Denotes if dependent services should also be deduplicated",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"time_window": {
				Description: "integer for time_unit",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
			},
			"time_unit": {
				Description: "time unit (mins or hours)",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "hour",
			},
			"basic_expressions": {
				Description: "The basic expression which needs to be evaluated to be true for this rule to apply.",
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lhs": {
							Description: "left hand side dropdown value",
							Type:        schema.TypeString,
							Required:    true,
						},
						"op": {
							Description:  "operator (is, is_not, matches, not_contains)",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"is", "is_not", "matches", "not_contains"}, false),
						},
						"rhs": {
							Description: "right hand side value",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceDeduplicationRuleImportV2(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	serviceID, ruleID, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("service_id", serviceID)
	d.SetId(ruleID)

	return []*schema.ResourceData{d}, nil
}

func resourceDeduplicationRuleCreateV2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	req := api.DeduplicationRule{
		IsBasic:                 d.Get("is_basic").(bool),
		Description:             d.Get("description").(string),
		Expression:              d.Get("expression").(string),
		DependencyDeduplication: d.Get("dependency_deduplication").(bool),
		TimeUnit:                d.Get("time_unit").(string),
		TimeWindow:              d.Get("time_window").(int),
	}

	basicExpressions, errx := decodeDeduplicationRuleBasicExpression(req.IsBasic, d.Get("basic_expressions").([]interface{}))
	if errx != nil {
		return errx
	}
	req.BasicExpression = basicExpressions

	tflog.Info(ctx, "Creating deduplication_rules", tf.M{
		"service_id": d.Get("service_id").(string),
	})

	dedupRule, err := client.CreateDeduplicationRulesV2(ctx, d.Get("service_id").(string), &api.CreateDeduplicationRule{Rule: req})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dedupRule.Rule.ID)

	return resourceDeduplicationRuleReadV2(ctx, d, meta)
}

func resourceDeduplicationRuleReadV2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	serviceID, ok := d.GetOk("service_id")
	if !ok {
		return diag.Errorf("invalid service id provided")
	}

	tflog.Info(ctx, "Reading deduplication_rules", tf.M{
		"id":         d.Id(),
		"service_id": d.Get("service_id").(string),
	})

	deduplicationRules, err := client.GetDeduplicationRuleByID(ctx, serviceID.(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(deduplicationRules, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDeduplicationRuleUpdateV2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	req := api.DeduplicationRule{
		IsBasic:                 d.Get("is_basic").(bool),
		Description:             d.Get("description").(string),
		Expression:              d.Get("expression").(string),
		DependencyDeduplication: d.Get("dependency_deduplication").(bool),
		TimeUnit:                d.Get("time_unit").(string),
		TimeWindow:              d.Get("time_window").(int),
	}

	basicExpressions, errx := decodeDeduplicationRuleBasicExpression(req.IsBasic, d.Get("basic_expressions").([]interface{}))
	if errx != nil {
		return errx
	}
	req.BasicExpression = basicExpressions

	_, err := client.UpdateDeduplicationRuleByID(ctx, d.Get("service_id").(string), d.Id(), &api.CreateDeduplicationRule{Rule: req})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDeduplicationRuleReadV2(ctx, d, meta)
}

func resourceDeduplicationRuleDeleteV2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteDeduplicationRuleByID(ctx, d.Get("service_id").(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func decodeDeduplicationRuleBasicExpression(isBasic bool, basicExpressions []interface{}) ([]*api.DeduplicationRuleCondition, diag.Diagnostics) {
	basicExpressionsReq := []*api.DeduplicationRuleCondition{}
	if (!isBasic && len(basicExpressions) > 0) || (isBasic && len(basicExpressions) == 0) {
		return nil, diag.Errorf("basic_expressions should be provided when is_basic is set to true, and should not be provided otherwise")
	}
	for _, expr := range basicExpressions {
		basicExpression, ok := expr.(map[string]interface{})
		if !ok {
			return nil, diag.Errorf("invalid basic expression format")
		}
		basicExpressionsReq = append(basicExpressionsReq, &api.DeduplicationRuleCondition{
			LHS: basicExpression["lhs"].(string),
			Op:  basicExpression["op"].(string),
			RHS: basicExpression["rhs"].(string),
		})
	}
	return basicExpressionsReq, nil
}
