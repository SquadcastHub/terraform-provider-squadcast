package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

const deduplicationRulesID = "deduplication_rules"

func resourceDeduplicationRules() *schema.Resource {
	return &schema.Resource{
		Description: "[Deduplication rules](https://support.squadcast.com/docs/de-duplication-rules) can help you reduce alert noise by organising and grouping alerts. This also provides easy access to similar alerts when needed. When these rules evaluate to true for an incoming incident, alerts will get deduplicated.",

		CreateContext: resourceDeduplicationRulesCreate,
		ReadContext:   resourceDeduplicationRulesRead,
		UpdateContext: resourceDeduplicationRulesUpdate,
		DeleteContext: resourceDeduplicationRulesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDeduplicationRulesImport,
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
						"description": {
							Description: "description.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"expression": {
							Description: "expression.",
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
							Description: "time window.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
						},
						"time_unit": {
							Description: "time unit.",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "hour",
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
									"op": {
										Description: "op",
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

func resourceDeduplicationRulesImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	teamID, serviceID, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("team_id", teamID)
	d.Set("service_id", serviceID)
	d.SetId(deduplicationRulesID)

	return []*schema.ResourceData{d}, nil
}

func resourceDeduplicationRulesCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	var rules []api.DeduplicationRule
	err := Decode(d.Get("rules"), &rules)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "Creating deduplication_rules", tf.M{
		"team_id":    d.Get("team_id").(string),
		"service_id": d.Get("service_id").(string),
	})

	_, err = client.UpdateDeduplicationRules(ctx, d.Get("service_id").(string), d.Get("team_id").(string), &api.UpdateDeduplicationRulesReq{Rules: rules})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(deduplicationRulesID)

	return resourceDeduplicationRulesRead(ctx, d, meta)
}

func resourceDeduplicationRulesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	serviceID, ok := d.GetOk("service_id")
	if !ok {
		return diag.Errorf("invalid service id provided")
	}

	tflog.Info(ctx, "Reading deduplication_rules", tf.M{
		"id":         d.Id(),
		"team_id":    d.Get("team_id").(string),
		"service_id": d.Get("service_id").(string),
	})
	deduplicationRules, err := client.GetDeduplicationRules(ctx, serviceID.(string), teamID.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(deduplicationRules, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDeduplicationRulesUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	var rules []api.DeduplicationRule
	err := Decode(d.Get("rules"), &rules)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateDeduplicationRules(ctx, d.Get("service_id").(string), d.Get("team_id").(string), &api.UpdateDeduplicationRulesReq{Rules: rules})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDeduplicationRulesRead(ctx, d, meta)
}

func resourceDeduplicationRulesDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.UpdateDeduplicationRules(ctx, d.Get("service_id").(string), d.Get("team_id").(string), &api.UpdateDeduplicationRulesReq{Rules: []api.DeduplicationRule{}})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
