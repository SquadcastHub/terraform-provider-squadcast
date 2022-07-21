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

const routingRulesID = "routing_rules"

func resourceRoutingRules() *schema.Resource {
	return &schema.Resource{
		Description: "[Routing rules](https://support.squadcast.com/docs/alert-routing) allows you to ensure that alerts are routed to the right responder with the help of `event tags` attached to them.",

		CreateContext: resourceRoutingRulesCreate,
		ReadContext:   resourceRoutingRulesRead,
		UpdateContext: resourceRoutingRulesUpdate,
		DeleteContext: resourceRoutingRulesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRoutingRulesImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"team_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
				ForceNew:     true,
			},
			"service_id": {
				Description:  "Service id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
				ForceNew:     true,
			},
			"rules": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_basic": {
							Description: "is basic?.",
							Type:        schema.TypeBool,
							Required:    true,
						},
						"expression": {
							Description: "expression.",
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
							Description:  "Type of the entity for which we are routing this incident - User, Escalation Policy or Squad",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"user", "escalationpolicy", "squad"}, false),
						},
						"basic_expressions": {
							Description: "basic expression.",
							Type:        schema.TypeList,
							Optional:    true,
							MinItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lhs": {
										Description: "lhs",
										Type:        schema.TypeString,
										Required:    true,
									},
									"rhs": {
										Description: "rhs",
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceRoutingRulesImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	teamID, serviceID, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("team_id", teamID)
	d.Set("service_id", serviceID)
	d.SetId(routingRulesID)

	return []*schema.ResourceData{d}, nil
}

func resourceRoutingRulesCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	var rules []api.RoutingRule
	err := Decode(d.Get("rules"), &rules)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "Creating routing_rules", tf.M{
		"team_id":    d.Get("team_id").(string),
		"service_id": d.Get("service_id").(string),
	})

	_, err = client.UpdateRoutingRules(ctx, d.Get("service_id").(string), d.Get("team_id").(string), &api.UpdateRoutingRulesReq{Rules: rules})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(routingRulesID)

	return resourceRoutingRulesRead(ctx, d, meta)
}

func resourceRoutingRulesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	serviceID, ok := d.GetOk("service_id")
	if !ok {
		return diag.Errorf("invalid service id provided")
	}

	tflog.Info(ctx, "Reading routing_rules", tf.M{
		"id":         d.Id(),
		"team_id":    d.Get("team_id").(string),
		"service_id": d.Get("service_id").(string),
	})
	routingRules, err := client.GetRoutingRules(ctx, serviceID.(string), teamID.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(routingRules, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRoutingRulesUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	var rules []api.RoutingRule
	err := Decode(d.Get("rules"), &rules)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateRoutingRules(ctx, d.Get("service_id").(string), d.Get("team_id").(string), &api.UpdateRoutingRulesReq{Rules: rules})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRoutingRulesRead(ctx, d, meta)
}

func resourceRoutingRulesDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.UpdateRoutingRules(ctx, d.Get("service_id").(string), d.Get("team_id").(string), &api.UpdateRoutingRulesReq{Rules: []api.RoutingRule{}})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
