package api

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type NewRotation struct {
	ID                          int                `graphql:"ID" json:"id,omitempty" tf:"id"`
	Name                        string             `graphql:"name" json:"name" tf:"name"`
	ParticipantGroups           []ParticipantGroup `graphql:"participantGroups" json:"participantGroups,omitempty" tf:"participant_groups"`
	StartDate                   string             `graphql:"startDate" json:"startDate" tf:"start_date"`
	Period                      string             `graphql:"period" json:"period" tf:"period"`
	ShiftTimeSlots              []Timeslot         `graphql:"shiftTimeSlots" json:"shiftTimeSlots" tf:"shift_timeslots"`
	CustomPeriodFrequency       int                `graphql:"customPeriodFrequency" json:"customPeriodFrequency,omitempty" tf:"custom_period_frequency"`
	CustomPeriodUnit            string             `graphql:"customPeriodUnit" json:"customPeriodUnit,omitempty" tf:"custom_period_unit"`
	ChangeParticipantsFrequency int                `graphql:"changeParticipantsFrequency" json:"changeParticipantsFrequency" tf:"change_participants_frequency"`
	ChangeParticipantsUnit      string             `graphql:"changeParticipantsUnit" json:"changeParticipantsUnit" tf:"change_participants_unit"`
	EndDate                     string             `graphql:"endDate" json:"endDate,omitempty" tf:"end_date"`
	EndsAfterIterations         int                `graphql:"endsAfterIterations" json:"endsAfterIterations,omitempty" tf:"ends_after_iterations"`
}

type ParticipantGroup struct {
	Participants []Participant `graphql:"participants" json:"participants" tf:"participants"`
}

type Participant struct {
	ID   string `graphql:"ID" json:"ID" tf:"id"`
	Type string `graphql:"type" json:"type" tf:"type"`
}

type Timeslot struct {
	StartHour   int    `graphql:"startHour" json:"startHour" tf:"start_hour"`
	StartMinute int    `graphql:"startMin" json:"startMin" tf:"start_minute"`
	Duration    int    `graphql:"duration" json:"duration" tf:"duration"`
	DayOfWeek   string `graphql:"dayOfWeek" json:"dayOfWeek,omitempty" tf:"day_of_week"`
}

// GraphQL query structs
type ScheduleRotationQueryStruct struct {
	NewRotation `graphql:"rotation(ID: $ID)"`
}

type ScheduleRotationMutateStruct struct {
	NewRotation `graphql:"createRotation(scheduleID: $scheduleID, input: $input)"`
}

type ScheduleRotationMutateDeleteStruct struct {
	NewRotation `graphql:"deleteRotation(ID: $ID)"`
}

func (ts Timeslot) Encode() (tf.M, error) {
	return tf.Encode(ts)
}

func (pg ParticipantGroup) Encode() (tf.M, error) {
	m, err := tf.Encode(pg)
	if err != nil {
		return nil, err
	}
	participantEncoded, perr := tf.EncodeSlice(pg.Participants)
	if perr != nil {
		return nil, perr
	}
	m["participants"] = participantEncoded
	return m, nil
}

func (p Participant) Encode() (tf.M, error) {
	return tf.Encode(p)
}

func (rot NewRotation) Encode() (tf.M, error) {
	m, err := tf.Encode(rot)
	if err != nil {
		return nil, err
	}

	timeslotsEncoded, terr := tf.EncodeSlice(rot.ShiftTimeSlots)
	if terr != nil {
		return nil, terr
	}
	m["shift_timeslots"] = timeslotsEncoded

	if rot.ParticipantGroups != nil {
		participantGroupsEncoded, perr := tf.EncodeSlice(rot.ParticipantGroups)
		if perr != nil {
			return nil, perr
		}
		m["participant_groups"] = participantGroupsEncoded
	}

	return m, nil
}

// ScheduleV2 APIs
func (client *Client) DeleteScheduleRotationByID(ctx context.Context, ID string) (*ScheduleRotationMutateDeleteStruct, error) {
	var m ScheduleRotationMutateDeleteStruct

	id, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		diag.Errorf("unable to convert schedule ID to string")
	}

	variables := map[string]interface{}{
		"ID": id,
	}

	return GraphQLRequest[ScheduleRotationMutateDeleteStruct]("mutate", client, ctx, &m, variables)
}

func (client *Client) GetScheduleRotationById(ctx context.Context, ID string) (*ScheduleRotationQueryStruct, error) {
	var m ScheduleRotationQueryStruct

	id, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		diag.Errorf("unable to convert schedule ID to string")
	}

	variables := map[string]interface{}{
		"ID": id,
	}

	return GraphQLRequest[ScheduleRotationQueryStruct]("query", client, ctx, &m, variables)
}

func (client *Client) CreateScheduleRotation(ctx context.Context, scheduleID int, payload NewRotation) (*ScheduleRotationMutateStruct, error) {
	var m ScheduleRotationMutateStruct

	variables := map[string]interface{}{
		"input":      payload,
		"scheduleID": scheduleID,
	}

	return GraphQLRequest[ScheduleRotationMutateStruct]("mutate", client, ctx, &m, variables)
}
