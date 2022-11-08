package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

const taggingRulesID = "tagging_rules"

func resourceTaggingRules() *schema.Resource {
	return &schema.Resource{
		Description: "[Tagging](https://support.squadcast.com/docs/event-tagging) is a rule-based, auto-tagging system with which you can define customised tags based on incident payloads, that get automatically assigned to incidents when they are triggered.",

		CreateContext: resourceTaggingRulesCreate,
		ReadContext:   resourceTaggingRulesRead,
		UpdateContext: resourceTaggingRulesUpdate,
		DeleteContext: resourceTaggingRulesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTaggingRulesImport,
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
										Description: "operator",
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
				},
			},
		},
	}
}

func resourceTaggingRulesImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	teamID, serviceID, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("team_id", teamID)
	d.Set("service_id", serviceID)
	d.SetId(taggingRulesID)

	return []*schema.ResourceData{d}, nil
}

// func decodeTaggingRules(input []any, output *[]api.TaggingRule) error {}

func resourceTaggingRulesCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	mrules := d.Get("rules").([]any)
	var rules []api.TaggingRule
	err := Decode(mrules, &rules)
	if err != nil {
		return diag.FromErr(err)
	}

	for i, mrule := range mrules {

		mtags := mrule.(tf.M)["tags"].([]any)

		tags := make(map[string]api.TaggingRuleTagValue, len(mtags))

		for _, mtag := range mtags {
			var tagvalue api.TaggingRuleTagValue
			err := Decode(mtag, &tagvalue)
			if err != nil {
				return diag.FromErr(err)
			}

			key := mtag.(tf.M)["key"].(string)

			tags[key] = tagvalue
		}

		rules[i].Tags = tags
	}

	tflog.Info(ctx, "Creating tagging_rules", tf.M{
		"team_id":    d.Get("team_id").(string),
		"service_id": d.Get("service_id").(string),
	})

	_, err = client.UpdateTaggingRules(ctx, d.Get("service_id").(string), d.Get("team_id").(string), &api.UpdateTaggingRulesReq{Rules: rules})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(taggingRulesID)

	return resourceTaggingRulesRead(ctx, d, meta)
}

func resourceTaggingRulesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	serviceID, ok := d.GetOk("service_id")
	if !ok {
		return diag.Errorf("invalid service id provided")
	}

	tflog.Info(ctx, "Reading tagging_rules", tf.M{
		"id":         d.Id(),
		"team_id":    d.Get("team_id").(string),
		"service_id": d.Get("service_id").(string),
	})
	taggingRules, err := client.GetTaggingRules(ctx, serviceID.(string), teamID.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(taggingRules, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceTaggingRulesUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	mrules := d.Get("rules").([]any)
	var rules []api.TaggingRule
	err := Decode(mrules, &rules)
	if err != nil {
		return diag.FromErr(err)
	}

	for i, mrule := range mrules {

		mtags := mrule.(tf.M)["tags"].([]any)

		tags := make(map[string]api.TaggingRuleTagValue, len(mtags))

		for _, mtag := range mtags {
			var tagvalue api.TaggingRuleTagValue
			err := Decode(mtag, &tagvalue)
			if err != nil {
				return diag.FromErr(err)
			}

			key := mtag.(tf.M)["key"].(string)

			tags[key] = tagvalue
		}

		rules[i].Tags = tags
	}

	_, err = client.UpdateTaggingRules(ctx, d.Get("service_id").(string), d.Get("team_id").(string), &api.UpdateTaggingRulesReq{Rules: rules})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTaggingRulesRead(ctx, d, meta)
}

func resourceTaggingRulesDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.UpdateTaggingRules(ctx, d.Get("service_id").(string), d.Get("team_id").(string), &api.UpdateTaggingRulesReq{Rules: []api.TaggingRule{}})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
