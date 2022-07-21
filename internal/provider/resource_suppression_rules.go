package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

const suppressionRulesID = "suppression_rules"

func resourceSuppressionRules() *schema.Resource {
	return &schema.Resource{
		Description: "[Suppression rules](https://support.squadcast.com/docs/alert-suppression) can help you avoid alert fatigue by suppressing notifications for non-actionable alerts." +

			"Squadcast will suppress the incidents that match any of the Suppression Rules you create for your Services. These incidents will go into the Suppressed state and you will not get any notifications for them",

		CreateContext: resourceSuppressionRulesCreate,
		ReadContext:   resourceSuppressionRulesRead,
		UpdateContext: resourceSuppressionRulesUpdate,
		DeleteContext: resourceSuppressionRulesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSuppressionRulesImport,
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

func resourceSuppressionRulesImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	teamID, serviceID, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("team_id", teamID)
	d.Set("service_id", serviceID)
	d.SetId(suppressionRulesID)

	return []*schema.ResourceData{d}, nil
}

func Decode(input any, output any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:               output,
		TagName:              tf.EncoderStructTag,
		ZeroFields:           true,
		IgnoreUntaggedFields: true,
	})
	if err != nil {
		return err
	}

	err = decoder.Decode(input)
	if err != nil {
		return err
	}

	return nil
}

func resourceSuppressionRulesCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	var rules []api.SuppressionRule
	err := Decode(d.Get("rules"), &rules)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "Creating suppression_rules", tf.M{
		"team_id":    d.Get("team_id").(string),
		"service_id": d.Get("service_id").(string),
	})

	_, err = client.UpdateSuppressionRules(ctx, d.Get("service_id").(string), d.Get("team_id").(string), &api.UpdateSuppressionRulesReq{Rules: rules})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(suppressionRulesID)

	return resourceSuppressionRulesRead(ctx, d, meta)
}

func resourceSuppressionRulesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	serviceID, ok := d.GetOk("service_id")
	if !ok {
		return diag.Errorf("invalid service id provided")
	}

	tflog.Info(ctx, "Reading suppression_rules", tf.M{
		"id":         d.Id(),
		"team_id":    d.Get("team_id").(string),
		"service_id": d.Get("service_id").(string),
	})
	suppressionRules, err := client.GetSuppressionRules(ctx, serviceID.(string), teamID.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(suppressionRules, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSuppressionRulesUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	var rules []api.SuppressionRule
	err := Decode(d.Get("rules"), &rules)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateSuppressionRules(ctx, d.Get("service_id").(string), d.Get("team_id").(string), &api.UpdateSuppressionRulesReq{Rules: rules})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSuppressionRulesRead(ctx, d, meta)
}

func resourceSuppressionRulesDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.UpdateSuppressionRules(ctx, d.Get("service_id").(string), d.Get("team_id").(string), &api.UpdateSuppressionRulesReq{Rules: []api.SuppressionRule{}})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
