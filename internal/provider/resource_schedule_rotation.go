package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func resourceScheduleRotation() *schema.Resource {
	return &schema.Resource{
		Description:   "[Schedule rotations](https://support.squadcast.com/schedules/schedules-new/adding-a-schedule#2.-choose-a-rotation-pattern) are used to manage on-call scheduling & determine who will be notified when an incident is triggered.",
		ReadContext:   resourceScheduleRotationRead,
		CreateContext: resourceScheduleRotationCreate,
		UpdateContext: resourceScheduleRotationUpdate,
		DeleteContext: resourceScheduleRotationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceScheduleRotationImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Rotation id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"schedule_id": {
				Description: "id of the schedule that the rotation belongs to.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"name": {
				Description:  "Rotation name.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 150),
			},
			"participant_groups": {
				Description: "Ordered list of participant groups for the rotation. For each rotation the participant_groups are cycled through in order.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"participants": {
							Description: "Group participants.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Description:  "Participant type (user, team, squad).",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"user", "squad", "team"}, false),
									},
									"id": {
										Description:  "Participant id.",
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
				Description: "Defines the start date of the rotation.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"period": {
				Description:  "Rotation period (none, daily, weekly, monthly, custom). Defines how often the rotation repeats.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "daily", "weekly", "monthly", "custom"}, false),
			},
			"shift_timeslots": {
				Description: "Timeslots where the rotation is active.",
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_hour": {
							Description:  "Defines the start hour of the each shift in the schedule timezone.",
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 23),
						},
						"start_minute": {
							Description:  "Defines the start minute of the each shift in the schedule timezone.",
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 59),
						},
						"duration": {
							Description:  "Defines the duration of each shift. (in minutes)",
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 1440),
						},
						"day_of_week": {
							Description:  "Defines the day of the week for the shift. If not specified, the timeslot is active on all days of the week.",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}, false),
						},
					},
				},
			},
			"custom_period_frequency": {
				Description:  "Frequency of the custom rotation repeat pattern. Only applicable if period is set to custom.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"custom_period_unit": {
				Description:  "Unit of the custom rotation repeat pattern (day, week, month). Only applicable if period is set to custom.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"day", "week", "month"}, false),
			},
			"change_participants_frequency": {
				Description:  "Frequency with which participants change in the rotation.",
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"change_participants_unit": {
				Description:  "Unit of the frequency with which participants change in the rotation (rotation, day, week, month).",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"rotation", "day", "week", "month"}, false),
			},
			"end_date": {
				Description: "Defines the end date of the schedule rotation.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ends_after_iterations": {
				Description: "Defines the number of iterations of the schedule rotation.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
		},
	}
}
func parse3PartImportID(id string) (string, string, string, error) {
	parts := strings.SplitN(id, ":", 3)

	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("unexpected format of import resource id (%s), expected teamID:ID", id)
	}

	return parts[0], parts[1], parts[2], nil
}

func resourceScheduleRotationImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*api.Client)
	teamID, scheduleName, rotationName, err := parse3PartImportID(d.Id())
	if err != nil {
		return nil, err
	}

	rotation, err := client.GetRotationByName(ctx, teamID, scheduleName, rotationName)
	if err != nil {
		return nil, errors.New("rotation not found")
	}
	d.SetId(strconv.Itoa(rotation.NewRotation.ID))

	return []*schema.ResourceData{d}, nil
}

func resourceScheduleRotationRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	id := d.Id()
	tflog.Info(ctx, "Reading rotation", tf.M{
		"id":   d.Id(),
		"name": d.Get("name").(string),
	})

	rotation, err := client.GetScheduleRotationById(ctx, id)
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(rotation, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceScheduleRotationCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Creating rotation", tf.M{
		"name": d.Get("name").(string),
	})

	createScheduleRotationReq := api.NewRotation{
		Name:                        d.Get("name").(string),
		StartDate:                   d.Get("start_date").(string),
		Period:                      d.Get("period").(string),
		ChangeParticipantsFrequency: d.Get("change_participants_frequency").(int),
		ChangeParticipantsUnit:      d.Get("change_participants_unit").(string),
	}

	endsAfterIterations, isIterationsEndSet := d.GetOk("ends_after_iterations")
	endDate, isEndDateSet := d.GetOk("end_date")
	if isIterationsEndSet {
		createScheduleRotationReq.EndsAfterIterations = endsAfterIterations.(int)
	}
	if isEndDateSet {
		createScheduleRotationReq.EndDate = endDate.(string)
	}
	if isIterationsEndSet && isEndDateSet {
		return diag.Errorf("only one of end_date and ends_after_iterations can be set")
	}

	participants := d.Get("participant_groups").([]interface{})
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
		createScheduleRotationReq.ParticipantGroups = participantGroupsList
	}

	shiftTimeSlots := d.Get("shift_timeslots").([]interface{})
	if len(shiftTimeSlots) > 0 {
		if createScheduleRotationReq.Period != "custom" && len(shiftTimeSlots) > 1 {
			return diag.Errorf("multiple shift_timeslots can only be set when period is custom")
		}
		var shiftTimeSlotsList []api.Timeslot
		err := Decode(shiftTimeSlots, &shiftTimeSlotsList)
		if err != nil {
			return diag.Errorf("shift_timeslots is invalid")
		}
		createScheduleRotationReq.ShiftTimeSlots = shiftTimeSlotsList
	}

	customPeriodFreq, freqIsSet := d.GetOk("custom_period_frequency")
	customPeriodUnit, unitIsSet := d.GetOk("custom_period_unit")

	if createScheduleRotationReq.Period == "custom" {
		if !freqIsSet || customPeriodFreq.(int) == 0 {
			return diag.Errorf("custom_period_frequency must be set when period is custom")
		}
		if !unitIsSet || customPeriodUnit.(string) == "" {
			return diag.Errorf("custom_period_unit must be set when period is custom")
		}

		createScheduleRotationReq.CustomPeriodFrequency = customPeriodFreq.(int)
		createScheduleRotationReq.CustomPeriodUnit = customPeriodUnit.(string)
	} else {
		if freqIsSet {
			return diag.Errorf("custom_period_frequency can only be set when period is custom")
		}
		if unitIsSet {
			return diag.Errorf("custom_period_unit can only be set when period is custom")
		}
	}

	rotation, err := client.CreateScheduleRotation(ctx, d.Get("schedule_id").(int), createScheduleRotationReq)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(rotation.NewRotation.ID))

	return resourceScheduleRotationRead(ctx, d, meta)
}

func resourceScheduleRotationUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Creating rotation", tf.M{
		"name": d.Get("name").(string),
	})
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	updateScheduleRotationReq := api.NewRotation{
		Name:                        d.Get("name").(string),
		StartDate:                   d.Get("start_date").(string),
		Period:                      d.Get("period").(string),
		ChangeParticipantsFrequency: d.Get("change_participants_frequency").(int),
		ChangeParticipantsUnit:      d.Get("change_participants_unit").(string),
	}

	endsAfterIterations, isIterationsEndSet := d.GetOk("ends_after_iterations")
	endDate, isEndDateSet := d.GetOk("end_date")
	if isIterationsEndSet {
		updateScheduleRotationReq.EndsAfterIterations = endsAfterIterations.(int)
	}
	if isEndDateSet {
		updateScheduleRotationReq.EndDate = endDate.(string)
	}
	if isIterationsEndSet && isEndDateSet {
		return diag.Errorf("only one of end_date and ends_after_iterations can be set")
	}

	participants := d.Get("participant_groups").([]interface{})
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
		updateScheduleRotationReq.ParticipantGroups = participantGroupsList
	}

	shiftTimeSlots := d.Get("shift_timeslots").([]interface{})
	if len(shiftTimeSlots) > 0 {
		if updateScheduleRotationReq.Period != "custom" && len(shiftTimeSlots) > 1 {
			return diag.Errorf("multiple shift_timeslots can only be set when period is custom")
		}
		var shiftTimeSlotsList []api.Timeslot
		err := Decode(shiftTimeSlots, &shiftTimeSlotsList)
		if err != nil {
			return diag.Errorf("shift_timeslots is invalid")
		}
		updateScheduleRotationReq.ShiftTimeSlots = shiftTimeSlotsList
	}

	customPeriodFreq, freqIsSet := d.GetOk("custom_period_frequency")
	customPeriodUnit, unitIsSet := d.GetOk("custom_period_unit")

	if updateScheduleRotationReq.Period == "custom" {
		if !freqIsSet || customPeriodFreq.(int) == 0 {
			return diag.Errorf("custom_period_frequency must be set when period is custom")
		}
		if !unitIsSet || customPeriodUnit.(string) == "" {
			return diag.Errorf("custom_period_unit must be set when period is custom")
		}

		updateScheduleRotationReq.CustomPeriodFrequency = customPeriodFreq.(int)
		updateScheduleRotationReq.CustomPeriodUnit = customPeriodUnit.(string)
	} else {
		if freqIsSet {
			return diag.Errorf("custom_period_frequency can only be set when period is custom")
		}
		if unitIsSet {
			return diag.Errorf("custom_period_unit can only be set when period is custom")
		}
	}

	_, err = client.UpdateScheduleRotation(ctx, id, updateScheduleRotationReq)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScheduleRotationRead(ctx, d, meta)
}

func resourceScheduleRotationDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteScheduleRotationByID(ctx, d.Id())
	if err != nil {
		tflog.Info(ctx, "No err while deleting rotation")
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			tflog.Info(ctx, "No resource found while deleting rotation")
			return nil
		}
		tflog.Info(ctx, "random err found while deleting rotation")
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "No err while deleting rotation")
	return nil
}
