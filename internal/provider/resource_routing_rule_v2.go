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

func resourceRoutingRuleV2() *schema.Resource {
	return &schema.Resource{
		Description: "[Routing rules](https://support.squadcast.com/docs/alert-routing) allows you to ensure that alerts are routed to the right responder with the help of `event tags` attached to them.",

		CreateContext: resourceRoutingRuleCreateV2,
		ReadContext:   resourceRoutingRuleReadV2,
		UpdateContext: resourceRoutingRuleUpdateV2,
		DeleteContext: resourceRoutingRuleDeleteV2,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRoutingRuleImportV2,
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
			"route_to_id": {
				Description:  "The id of the entity (user, escalation policy, squad) for which we are routing this incident.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
			},
			"route_to_type": {
				Description:  "Type of the entity for which we are routing this incident (user, escalationpolicy or squad)",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"user", "escalationpolicy", "squad"}, false),
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

func resourceRoutingRuleImportV2(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	serviceID, ruleID, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}
	d.Set("service_id", serviceID)
	d.SetId(ruleID)

	return []*schema.ResourceData{d}, nil
}

func resourceRoutingRuleCreateV2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	req := api.RoutingRule{
		IsBasic:    d.Get("is_basic").(bool),
		Expression: d.Get("expression").(string),
		RouteTo: api.RouteTo{
			EntityID:   d.Get("route_to_id").(string),
			EntityType: d.Get("route_to_type").(string),
		},
	}

	if req.IsBasic && len(req.Expression) > 0 {
		return diag.Errorf("expression should be passed only when is_basic is set to true")
	}

	basicExpressions, errx := decodeRoutingRuleBasicExpression(req.IsBasic, d.Get("basic_expressions").([]interface{}))
	if errx != nil {
		return errx
	}
	req.BasicExpression = basicExpressions

	tflog.Info(ctx, "Creating routing_rules", tf.M{
		"service_id": d.Get("service_id").(string),
	})

	routingRule, err := client.CreateRoutingRulesV2(ctx, d.Get("service_id").(string), &api.CreateRoutingRule{Rule: req})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(routingRule.Rule.ID)

	return resourceRoutingRuleReadV2(ctx, d, meta)
}

func resourceRoutingRuleReadV2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	serviceID, ok := d.GetOk("service_id")
	if !ok {
		return diag.Errorf("invalid service id provided")
	}

	tflog.Info(ctx, "Reading routing_rules", tf.M{
		"id":         d.Id(),
		"service_id": d.Get("service_id").(string),
	})
	routingRules, err := client.GetRoutingRuleByID(ctx, serviceID.(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(routingRules, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRoutingRuleUpdateV2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	req := api.RoutingRule{
		IsBasic:    d.Get("is_basic").(bool),
		Expression: d.Get("expression").(string),
		RouteTo: api.RouteTo{
			EntityID:   d.Get("route_to_id").(string),
			EntityType: d.Get("route_to_type").(string),
		},
		BasicExpression: []*api.RoutingRuleCondition{},
	}

	if req.IsBasic && len(req.Expression) > 0 {
		return diag.Errorf("expression should be passed only when is_basic is set to true")
	}

	basicExpressions, err := decodeRoutingRuleBasicExpression(req.IsBasic, d.Get("basic_expressions").([]interface{}))
	if err != nil {
		return err
	}
	req.BasicExpression = basicExpressions

	_, errx := client.UpdateRoutingRuleByID(ctx, d.Get("service_id").(string), d.Id(), &api.CreateRoutingRule{Rule: req})
	if errx != nil {
		return diag.FromErr(errx)
	}

	return resourceRoutingRuleReadV2(ctx, d, meta)
}

func resourceRoutingRuleDeleteV2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteRoutingRuleByID(ctx, d.Get("service_id").(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func decodeRoutingRuleBasicExpression(isBasic bool, basicExpressions []interface{}) ([]*api.RoutingRuleCondition, diag.Diagnostics) {
	basicExpressionsReq := []*api.RoutingRuleCondition{}
	if (!isBasic && len(basicExpressions) > 0) || (isBasic && len(basicExpressions) == 0) {
		return nil, diag.Errorf("basic_expressions should be provided when is_basic is set to true, and should not be provided otherwise")
	}
	for _, expr := range basicExpressions {
		basicExpression, ok := expr.(map[string]interface{})
		if !ok {
			return nil, diag.Errorf("invalid basic expression format")
		}
		basicExpressionsReq = append(basicExpressionsReq, &api.RoutingRuleCondition{
			LHS: basicExpression["lhs"].(string),
			RHS: basicExpression["rhs"].(string),
		})
	}
	return basicExpressionsReq, nil
}
