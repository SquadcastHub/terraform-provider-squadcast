package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
							Description: "is_basic will be true when users use the drop down selectors which will have lhs, op & rhs value, whereas it will be false when they use the advanced mode and it would have the expression for it's value",
							Type:        schema.TypeBool,
							Required:    true,
						},
						"is_timebased": {
							Description: "is_timebased will be true when users use the time based suppression rule",
							Type:        schema.TypeBool,
							Computed:    true,
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
						"timeslots": {
							Description: "The timeslots for which this rule should be applied.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time_zone": {
										Description: "Time zone for the time slot",
										Type:        schema.TypeString,
										Required:    true,
										ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
											_, err := time.LoadLocation(val.(string))
											if err != nil {
												errs = append(errs, err)
											}
											return
										},
									},
									"start_time": {
										Description: "Defines the start date of the time slot",
										Type:        schema.TypeString,
										Required:    true,
									},
									"end_time": {
										Description: "Defines the end date of the time slot",
										Type:        schema.TypeString,
										Required:    true,
									},
									"ends_on": {
										Description: "Defines the end date of the repetition",
										Type:        schema.TypeString,
										Required:    true,
									},
									"repetition": {
										Description:  "Defines the repetition of the time slot",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"none", "daily", "weekly", "monthly", "custom"}, false),
									},
									"is_allday": {
										Description: "Defines if the time slot is an all day slot",
										Type:        schema.TypeBool,
										Optional:    true,
									},
									"ends_never": {
										Description: "Defines whether the time slot ends or not",
										Type:        schema.TypeBool,
										Optional:    true,
									},
									"is_custom": {
										Description: "Defines whether repetition is custom or not",
										Type:        schema.TypeBool,
										Computed:    true,
									},
									"custom": {
										Description: "Use this field to specify the custom time slots for which this rule should be applied. This field is only applicable when the repetition field is set to custom.",
										Type:        schema.TypeList,
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"repeats": {
													Description:  "Determines how often the rule repeats. Valid values are day, week, month.",
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice([]string{"day", "week", "month"}, false),
												},
												"repeats_count": {
													Description: "Number of times to repeat.",
													Type:        schema.TypeInt,
													Optional:    true,
												},
												"repeats_on_month": {
													Description: "Repeats on month.",
													Type:        schema.TypeString,
													Computed:    true,
												},
												"repeats_on_weekdays": {
													Description: "List of weekdays to repeat on.",
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    7,
													Elem: &schema.Schema{
														Type:         schema.TypeInt,
														ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 3, 4, 5, 6}),
													},
												},
											},
										},
									},
								},
							},
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
	mrules := d.Get("rules").([]interface{})
	// if repetition is custom, convert timeslots.custom in each rule from list to map
	for i, rule := range mrules {
		mrule := rule.(map[string]interface{})
		mtimeSlots := mrule["timeslots"].([]interface{})
		if len(mtimeSlots) != 0 {
			for _, mtimeSlot := range mtimeSlots {
				mtimeSlot := mtimeSlot.(map[string]interface{})
				if mtimeSlot["repetition"] != "custom" { // if repetition is not custom, skip
					mtimeSlot["custom"] = nil
					continue
				}

				// else, get custom property and convert it to api.CustomTime
				/****************************************************
					tfstate format:
						"timeslots": [
							{
								....,
								"custom": [
									{
										...
									}
								]
							}
						]
				****************************************************/
				if len(mtimeSlot["custom"].([]interface{})) == 0 {
					return diag.Errorf("timeslot.custom is empty")
				}
				mcustom := mtimeSlot["custom"].([]interface{})[0].(map[string]interface{})
				mrepeats := mcustom["repeats"].(string)
				mrepeatOnWeekdays := mcustom["repeats_on_weekdays"].([]interface{})
				repeatOnWeekdays := make([]int, len(mrepeatOnWeekdays))
				repeatsOnMonth := ""

				// ? VALIDATION:
				// if repeats is week, set repeats_on_weekdays to the value from tfstate
				// if repeats is not week, set repeats_on_weekdays to nil
				// if repeats is month, set repeats_on_month to date-occurrence

				switch mrepeats {
				case "week":
					for i, v := range mrepeatOnWeekdays {
						repeatOnWeekdays[i] = v.(int)
					}
				case "month":
					repeatsOnMonth = "date-occurrence"
				default:
					repeatOnWeekdays = nil
				}
				// set custom property to api.CustomTime
				mtimeSlot["custom"] = api.CustomTime{
					RepeatsOnMonth:    repeatsOnMonth,
					RepeatsOnWeekdays: repeatOnWeekdays,
					RepeatsCount:      mcustom["repeats_count"].(int),
					Repeats:           mrepeats,
				}
				mtimeSlot["is_custom"] = true
			}
			// convert mtimeslots to api.TimeSlot
			var timeslots []*api.TimeSlot
			err := Decode(mtimeSlots, &timeslots)
			if err != nil {
				return diag.FromErr(err)
			}
			mrules[i].(map[string]interface{})["is_timebased"] = true
			mrules[i].(map[string]interface{})["timeslots"] = timeslots
		}
	}

	err := Decode(mrules, &rules)
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
	mrules := d.Get("rules").([]interface{})
	for i, rule := range mrules {
		mrule := rule.(map[string]interface{})
		mtimeSlots := mrule["timeslots"].([]interface{})
		if len(mtimeSlots) != 0 {
			for _, mtimeSlot := range mtimeSlots {
				mtimeSlot := mtimeSlot.(map[string]interface{})
				if mtimeSlot["repetition"] != "custom" {
					mtimeSlot["custom"] = nil
					continue
				}
				if len(mtimeSlot["custom"].([]interface{})) == 0 {
					return diag.Errorf("timeslots.custom cannot be empty when timeslots.repetition is set to 'custom'")
				}
				mcustom := mtimeSlot["custom"].([]interface{})[0].(map[string]interface{})
				mrepeats := mcustom["repeats"].(string)
				mrepeatOnWeekdays := mcustom["repeats_on_weekdays"].([]interface{})
				repeatOnWeekdays := make([]int, len(mrepeatOnWeekdays))
				repeatsOnMonth := ""

				switch mrepeats {
				case "week":
					for i, v := range mrepeatOnWeekdays {
						repeatOnWeekdays[i] = v.(int)
					}
				case "month":
					repeatsOnMonth = "date-occurrence"
				default:
					if len(mrepeatOnWeekdays) != 0 {
						return diag.Errorf("timeslots.custom.repeats_on_weekdays cannot be set when timeslots.custom.repeats is not set to 'week'")
					}
					repeatOnWeekdays = nil
				}

				mtimeSlot["custom"] = api.CustomTime{
					RepeatsOnMonth:    repeatsOnMonth,
					RepeatsOnWeekdays: repeatOnWeekdays,
					RepeatsCount:      mcustom["repeats_count"].(int),
					Repeats:           mrepeats,
				}
				mtimeSlot["is_custom"] = true
			}
			var timeslots []*api.TimeSlot
			err := Decode(mtimeSlots, &timeslots)
			if err != nil {
				return diag.FromErr(err)
			}
			mrules[i].(map[string]interface{})["is_timebased"] = true
			mrules[i].(map[string]interface{})["timeslots"] = timeslots
		}
	}

	err := Decode(mrules, &rules)
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
