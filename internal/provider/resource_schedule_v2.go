package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func resourceScheduleV2() *schema.Resource {
	return &schema.Resource{
		Description:   "Squadcast v2 schedules", //todo: update this
		ReadContext:   resourceScheduleV2Read,
		CreateContext: resourceScheduleV2Create,
		UpdateContext: resourceScheduleV2Create,
		DeleteContext: resourceScheduleV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceScheduleV2Import,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Schedule id.",
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
			"name": {
				Description:  "Name of the schedule.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": {
				Description:  "Detailed description about the Schedule.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"timezone": {
				Description: "Timezone of the schedule",
				Type:        schema.TypeString,
				Required:    true,
			},
			"entity_owner": {
				Description: "Schedule entity owner.",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description:  "Schedule owner type (user, team, squad).",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"user", "squad", "team"}, false),
						},
						"id": {
							Description:  "Schedule owner id.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tf.ValidateObjectID,
						},
					},
				},
			},
			"tags": {
				Description: "Schedule tags.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Description: "Schedule tag key.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"value": {
							Description: "Schedule tag value.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"color": {
							Description: "Schedule tag color.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"rotations": {
				Description: "Schedule rotations.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Schedule rotation id.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": {
							Description:  "Schedule rotation name.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"participant_groups": {
							Description: "Schedule rotation participant groups.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"participants": {
										Description: "Schedule rotation participants.",
										Type:        schema.TypeList,
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Description: "Schedule rotation participant type (user, team, squad).",
													Type:        schema.TypeString,
													Required:    true,
													ValidateFunc: validation.StringInSlice([]string{"user", "squad", "team"}, false),
												},
												"id": {
													Description:  "Schedule rotation participant id.",
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: tf.ValidateObjectID,
												},
											},
										},
									},
								},
							},
						},
						"start_date": {
							Description: "Schedule rotation start date.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"period": {
							Description:  "Schedule rotation period.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"none", "daily", "weekly", "monthly", "custom"}, false),
						},
						"shift_timeslots": {
							Description: "Schedule rotation shift timeslots.",
							Type:        schema.TypeList,
							Required:    true,
							MinItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_hour": {
										Description:  "Schedule rotation shift timeslots start hour.",
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 23),
									},
									"start_minute": {
										Description:  "Schedule rotation shift timeslots start minute.",
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 59),
									},
									"duration": {
										Description:  "Schedule rotation shift timeslots duration in minutes.",
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 1440),
									},
									"day_of_week": {
										Description:  "Schedule rotation shift timeslots day of week.",
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}, false),
									},
								},
							},
						},
						"custom_period_frequency": {
							Description: "Schedule rotation custom period frequency.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"custom_period_unit": {
							Description: "Schedule rotation custom period unit.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"change_participants_frequency": {
							Description: "Schedule rotation change participants frequency.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"change_participants_unit": {
							Description: "Schedule rotation change participants unit.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"end_date": {
							Description: "Schedule rotation end date.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"ends_after_iterations": {
							Description: "Schedule rotation ends after iterations.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourceScheduleV2Import(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*api.Client)

	schedule, err := client.GetScheduleV2ById(ctx, d.Id())
	if err != nil {
		return nil, err
	}
	d.SetId(strconv.Itoa(schedule.NewSchedule.ID))

	return []*schema.ResourceData{d}, nil
}

func resourceScheduleV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	id := d.Id()
	tflog.Info(ctx, "Reading schedule", tf.M{
		"id":   d.Id(),
		"name": d.Get("name").(string),
	})

	schedule, err := client.GetScheduleV2ById(ctx, id)
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(schedule, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceScheduleV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Creating schedule", tf.M{
		"name": d.Get("name").(string),
	})

	createScheduleReq := api.NewSchedule{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		TimeZone:    d.Get("timezone").(string),
		TeamID:      d.Get("team_id").(string),
	}

	tags := d.Get("tags").([]interface{})

	if len(tags) > 0 {
		var tagsList []*api.Tag
		err := Decode(tags, &tagsList)
		if err != nil {
			return diag.Errorf("tags is invalid")
		}
		createScheduleReq.Tags = tagsList
	}

	rotations := d.Get("rotations").([]interface{})
	if len(rotations) > 0 {
		var rotationsList []*api.Rotation
		for _, rotation := range rotations {
			rotationMap, ok := rotation.(map[string]interface{})
			if !ok {
				return diag.Errorf("rotation is invalid")
			}
			r := &api.Rotation{
				Name:                        rotationMap["name"].(string),
				Period:                      rotationMap["period"].(string),
				ChangeParticipantsFrequency: rotationMap["change_participants_frequency"].(int),
				ChangeParticipantsUnit:      rotationMap["change_participants_unit"].(string),
				StartDate:                   rotationMap["start_date"].(string),
				EndDate:                     rotationMap["end_date"].(string),
				EndsAfterIterations:         rotationMap["ends_after_iterations"].(int),
			}
			// convert participants to []api.Participant
			participants := rotationMap["participant_groups"].([]interface{})
			if len(participants) > 0 {
				var participantGroupsList []api.ParticipantGroup
				for _, participant := range participants {
					participantMap, ok := participant.(map[string]interface{})
					if !ok {
						return diag.Errorf("participant_groups is invalid")
					}
					var participantGroup api.ParticipantGroup
					var participantsList []api.Participant
					participants := participantMap["participants"].([]interface{})

					err := Decode(participants, &participantsList)
					if err != nil {
						return diag.Errorf(err.Error())
					}
					participantGroup.Participants = participantsList
					participantGroupsList = append(participantGroupsList, participantGroup)
				}
				r.ParticipantGroups = participantGroupsList
			}

			// convert shift_timeslots to []api.Timeslot
			shiftTimeSlots := rotationMap["shift_timeslots"].([]interface{})
			if len(shiftTimeSlots) > 0 {
				if r.Period != "custom" && len(shiftTimeSlots) > 1 {
					return diag.Errorf("shift_timeslots can only have one timeslot when period is not custom")
				}
				var shiftTimeSlotsList []api.Timeslot
				err := Decode(shiftTimeSlots, &shiftTimeSlotsList)
				if err != nil {
					return diag.Errorf("shift_timeslots is invalid")
				}
				r.ShiftTimeSlots = shiftTimeSlotsList
			}

			if r.Period == "custom" {
				if val, ok := rotationMap["custom_period_frequency"].(int); ok {
					r.CustomPeriodFrequency = val
				} else {
					return diag.Errorf("custom_period_frequency must be set when period is custom")
				}
				if val, ok := rotationMap["custom_period_unit"].(string); ok {
					r.CustomPeriodUnit = val
				} else {
					return diag.Errorf("custom_period_unit must be set when period is custom")
				}
			}

			rotationsList = append(rotationsList, r)
		}
		createScheduleReq.Rotations = rotationsList
	}

	entityOwner := d.Get("entity_owner").([]interface{})
	if len(entityOwner) > 0 {
		entityOwnerMap, ok := entityOwner[0].(map[string]interface{})
		if !ok {
			return diag.Errorf("entity_owner is invalid")
		}
		createScheduleReq.Owner = &api.Owner{
			Type: entityOwnerMap["type"].(string),
			ID:   entityOwnerMap["id"].(string),
		}
	}

	schedule, err := client.CreateScheduleV2(ctx, createScheduleReq)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(schedule.NewSchedule.ID))

	return resourceScheduleV2Read(ctx, d, meta)
}

func resourceScheduleV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteScheduleV2ByID(ctx, d.Id())
	if err != nil {
		tflog.Info(ctx, "No err while deleting schedule")
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			tflog.Info(ctx, "No resource found while deleting schedule")
			return nil
		}
		tflog.Info(ctx, "random err found while deleting schedule")
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "No err while deleting schedule")
	return nil
}
