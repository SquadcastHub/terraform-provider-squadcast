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

func resourceTaggingRuleV2() *schema.Resource {
	return &schema.Resource{
		Description: "[Tagging](https://support.squadcast.com/docs/event-tagging) is a rule-based, auto-tagging system with which you can define customised tags based on incident payloads, that get automatically assigned to incidents when they are triggered.",

		CreateContext: resourceTaggingRuleCreateV2,
		ReadContext:   resourceTaggingRuleReadV2,
		UpdateContext: resourceTaggingRuleUpdateV2,
		DeleteContext: resourceTaggingRuleDeleteV2,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTaggingRuleImportV2,
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
			"expression": {
				Description: "The expression which needs to be evaluated to be true for this rule to apply.",
				Type:        schema.TypeString,
				Optional:    true,
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
			"tags": {
				Description: "The tags supposed to be set for a given payload(incident), Expression must be set when tags are empty and must contain addTags parameters.",
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
						"color": {
							Description: "Tag color, hex values",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceTaggingRuleImportV2(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	serviceID, ruleID, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}
	d.Set("service_id", serviceID)
	d.SetId(ruleID)

	return []*schema.ResourceData{d}, nil
}

func resourceTaggingRuleCreateV2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)
	req := api.TaggingRule{
		IsBasic:    d.Get("is_basic").(bool),
		Expression: d.Get("expression").(string),
	}
	if req.IsBasic && len(req.Expression) > 0 {
		return diag.Errorf("expression should be passed only when is_basic is set to true")
	}
	basicExpressions, errx := decodeTaggingRuleBasicExpression(req.IsBasic, d.Get("basic_expressions").([]interface{}))
	if errx != nil {
		return errx
	}
	req.BasicExpression = basicExpressions

	tags, errx := decodeTaggingRuleTags(d.Get("tags").([]interface{}))
	if errx != nil {
		return errx
	}
	if len(tags) > 0 {
		req.Tags = tags
	}

	tflog.Info(ctx, "Creating tagging_rules", tf.M{
		"service_id": d.Get("service_id").(string),
	})

	taggingRule, err := client.CreateTaggingRulesV2(ctx, d.Get("service_id").(string), &api.CreateTaggingRule{
		Rule: req,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(taggingRule.Rule.ID)

	return resourceTaggingRuleReadV2(ctx, d, meta)
}

func resourceTaggingRuleReadV2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	serviceID, ok := d.GetOk("service_id")
	if !ok {
		return diag.Errorf("invalid service id provided")
	}

	tflog.Info(ctx, "Reading tagging_rules", tf.M{
		"id":         d.Id(),
		"service_id": d.Get("service_id").(string),
	})
	taggingRule, err := client.GetTaggingRuleByID(ctx, serviceID.(string), d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(taggingRule, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceTaggingRuleUpdateV2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	req := api.TaggingRule{
		IsBasic:         d.Get("is_basic").(bool),
		Expression:      d.Get("expression").(string),
		BasicExpression: []*api.TaggingRuleCondition{},
		Tags:            map[string]api.TaggingRuleTagValue{},
	}

	if req.IsBasic && len(req.Expression) > 0 {
		return diag.Errorf("expression should be passed only when is_basic is set to true")
	}
	basicExpressions, errx := decodeTaggingRuleBasicExpression(req.IsBasic, d.Get("basic_expressions").([]interface{}))
	if errx != nil {
		return errx
	}
	req.BasicExpression = basicExpressions

	tags, errx := decodeTaggingRuleTags(d.Get("tags").([]interface{}))
	if errx != nil {
		return errx
	}
	if len(tags) > 0 {
		req.Tags = tags
	}

	tflog.Info(ctx, "Updating tagging_rules", tf.M{
		"id":         d.Id(),
		"service_id": d.Get("service_id").(string),
	})

	_, err := client.UpdateTaggingRuleByID(ctx, d.Get("service_id").(string), d.Id(), &api.CreateTaggingRule{Rule: req})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTaggingRuleReadV2(ctx, d, meta)
}

func resourceTaggingRuleDeleteV2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteTaggingRuleByID(ctx, d.Get("service_id").(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func decodeTaggingRuleBasicExpression(isBasic bool, basicExpressions []interface{}) ([]*api.TaggingRuleCondition, diag.Diagnostics) {
	basicExpressionsReq := []*api.TaggingRuleCondition{}
	if (!isBasic && len(basicExpressions) > 0) || (isBasic && len(basicExpressions) == 0) {
		return nil, diag.Errorf("basic_expressions should be provided when is_basic is set to true, and should not be provided otherwise")
	}

	for _, expr := range basicExpressions {
		basicExpression, ok := expr.(map[string]interface{})
		if !ok {
			return nil, diag.Errorf("invalid basic expression format")
		}
		basicExpressionsReq = append(basicExpressionsReq, &api.TaggingRuleCondition{
			LHS: basicExpression["lhs"].(string),
			Op:  basicExpression["op"].(string),
			RHS: basicExpression["rhs"].(string),
		})
	}
	return basicExpressionsReq, nil
}

func decodeTaggingRuleTags(mtags []interface{}) (map[string]api.TaggingRuleTagValue, diag.Diagnostics) {
	tags := make(map[string]api.TaggingRuleTagValue, len(mtags))

	if len(mtags) > 0 {
		for _, mtag := range mtags {
			var tagvalue api.TaggingRuleTagValue
			err := Decode(mtag, &tagvalue)
			if err != nil {
				return nil, diag.FromErr(err)
			}
			key := mtag.(tf.M)["key"].(string)
			tags[key] = tagvalue
		}

		return tags, nil
	}

	return tags, nil
}
