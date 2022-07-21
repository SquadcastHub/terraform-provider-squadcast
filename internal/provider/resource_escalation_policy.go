package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func resourceEscalationPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "[Escalation Policies](https://support.squadcast.com/docs/escalation-policies) defines rules indicating when and how alerts will escalate to various Users, Squads and (or) Schedules within your Organization.",

		CreateContext: resourceEscalationPolicyCreate,
		ReadContext:   resourceEscalationPolicyRead,
		UpdateContext: resourceEscalationPolicyUpdate,
		DeleteContext: resourceEscalationPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceEscalationPolicyImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "EscalationPolicy id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:  "Name of the Escalation Policy.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": {
				Description:  "Detailed description about the Escalation Policy.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"team_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
				ForceNew:     true,
			},
			"repeat": {
				Description: "You can choose to repeate the entire policy, if no one acknowledges the incident even after the Escalation Policy has been executed fully once",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"times": {
							Description: "The number of times you want this escalation policy to be repeated, maximum allowed to repeat 3 times",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"delay_minutes": {
							Description: "The number of minutes to wait before repeating the escalation policy",
							Type:        schema.TypeInt,
							Required:    true,
						},
					},
				},
			},
			"rules": {
				Description: "Rules will have the details of who to notify and when to notify and how to notify them.",
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delay_minutes": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"targets": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: tf.ValidateObjectID,
									},
									"type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"user", "squad", "schedule"}, false),
									},
								},
							},
						},
						"notification_channels": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"SMS", "Phone", "Email", "Push"}, false),
							},
						},
						"round_robin": {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Description: "Enables Round Robin escalation within this layer",
										Type:        schema.TypeBool,
										Required:    true,
									},
									"rotation": {
										Type:     schema.TypeList,
										Optional: true,
										MinItems: 1,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Description: "enable rotation within",
													Type:        schema.TypeBool,
													Optional:    true,
												},
												"delay_minutes": {
													Description: "repeat after minutes",
													Type:        schema.TypeInt,
													Optional:    true,
												},
											},
										},
									},
								},
							},
						},
						"repeat": {
							Description: "repeat this rule",
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"times": {
										Description: "repeat times",
										Type:        schema.TypeInt,
										Required:    true,
									},
									"delay_minutes": {
										Description: "repeat after minutes",
										Type:        schema.TypeInt,
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

func resourceEscalationPolicyImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*api.Client)

	teamID, name, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}

	escalationPolicy, err := client.GetEscalationPolicyByName(ctx, teamID, name)
	if err != nil {
		return nil, err
	}

	d.Set("team_id", teamID)
	d.SetId(escalationPolicy.ID)

	return []*schema.ResourceData{d}, nil
}

func decodeEscalationPolicyRules(mrules []tf.M) ([]api.EscalationPolicyRule, error) {
	rules := make([]api.EscalationPolicyRule, 0)

	for i, mrule := range mrules {
		rule := api.EscalationPolicyRule{
			EscalateAfterMinutes: mrule["delay_minutes"].(int),
		}

		mrepeats := tf.ListToSlice[tf.M](mrule["repeat"])
		if len(mrepeats) == 1 {
			rule.RepeatTimes = mrepeats[0]["times"].(int)
			rule.RepeatAfterMinutes = mrepeats[0]["delay_minutes"].(int)
		}

		mrr := tf.ListToSlice[tf.M](mrule["round_robin"])
		if len(mrr) == 1 {
			rule.RoundrobinEnabled = mrr[0]["enabled"].(bool)

			mrrrotation := tf.ListToSlice[tf.M](mrr[0]["rotation"])
			if len(mrrrotation) == 1 {
				rule.EscalateWithinRoundrobin = mrrrotation[0]["enabled"].(bool)
				rule.RepeatAfterMinutes = mrrrotation[0]["delay_minutes"].(int)
			}
		}

		if len(mrepeats) == 1 && rule.RoundrobinEnabled {
			return nil, fmt.Errorf("rule %d cannot have both round robin and a repetition, please remove one", i)
		}

		mtargets := tf.ListToSlice[tf.M](mrule["targets"])
		targets := make([]*api.EscalationPolicyTarget, 0)
		for _, mtarget := range mtargets {
			target := api.EscalationPolicyTarget{
				ID:   mtarget["id"].(string),
				Type: mtarget["type"].(string),
			}
			targets = append(targets, &target)
		}
		rule.Targets = targets

		rule.Via = tf.ListToSlice[string](mrule["notification_channels"])

		rules = append(rules, rule)
	}

	return rules, nil
}

func decodeEscalationPolicy(d *schema.ResourceData) (*api.CreateUpdateEscalationPolicyReq, error) {
	rules, err := decodeEscalationPolicyRules(tf.ListToSlice[tf.M](d.Get("rules")))
	if err != nil {
		return nil, fmt.Errorf("escalation policy `%s` is invalid: %s", d.Get("name").(string), err.Error())
	}

	req := &api.CreateUpdateEscalationPolicyReq{
		TeamID:             d.Get("team_id").(string),
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		RepeatTimes:        d.Get("repeat.0.times").(int),
		RepeatAfterMinutes: d.Get("repeat.0.delay_minutes").(int),
		Rules:              rules,
		IsUsingNewFields:   true,
	}

	return req, nil
}

func resourceEscalationPolicyCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Creating escalation_policy", tf.M{
		"name": d.Get("name").(string),
	})

	req, err := decodeEscalationPolicy(d)
	if err != nil {
		return diag.FromErr(err)
	}

	escalationPolicy, err := client.CreateEscalationPolicy(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(escalationPolicy.ID)

	return resourceEscalationPolicyRead(ctx, d, meta)
}

func resourceEscalationPolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	id := d.Id()

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	tflog.Info(ctx, "Reading escalation_policy", tf.M{
		"id":   d.Id(),
		"name": d.Get("name").(string),
	})
	escalationPolicy, err := client.GetEscalationPolicyById(ctx, teamID.(string), id)
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(escalationPolicy, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceEscalationPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	req, err := decodeEscalationPolicy(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateEscalationPolicy(ctx, d.Id(), req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceEscalationPolicyRead(ctx, d, meta)
}

func resourceEscalationPolicyDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteEscalationPolicy(ctx, d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
