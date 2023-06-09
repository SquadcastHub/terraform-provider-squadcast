package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

// legacy schedule
type Schedule struct {
	ID          string   `json:"id" tf:"id"`
	Name        string   `json:"name" tf:"name"`
	Slug        string   `json:"slug" tf:"-"`
	Colour      string   `json:"colour" tf:"color"`
	Description string   `json:"description" tf:"description"`
	Owner       OwnerRef `json:"owner" tf:"-"`
}

type NewSchedule struct {
	ID          int         `graphql:"ID" json:"id,omitempty" tf:"id"`
	Name        string      `graphql:"name" json:"name" tf:"name"`
	Description string      `graphql:"description" json:"description" tf:"description"`
	TimeZone    string      `graphql:"timeZone" json:"timeZone" tf:"timezone"`
	TeamID      string      `graphql:"teamID" json:"teamID" tf:"team_id"`
	Owner       *Owner      `graphql:"owner" json:"owner" tf:"-"`
	Tags        []*Tag      `graphql:"tags" json:"tags,omitempty" tf:"tags"`
	Rotations   []*Rotation `graphql:"rotations" json:"rotations,omitempty" tf:"rotations"`
}

type Owner struct {
	ID   string `graphql:"ID" json:"ID" tf:"id"`
	Type string `graphql:"type" json:"type" tf:"type"`
}

type Tag struct {
	Key   string `graphql:"key" json:"key" tf:"key"`
	Value string `graphql:"value" json:"value" tf:"value"`
	Color string `graphql:"color" json:"color" tf:"color"`
}

type Rotation struct {
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
type ScheduleQueryStruct struct {
	NewSchedule `graphql:"schedule(ID: $ID)"`
}

type ScheduleMutateStruct struct {
	NewSchedule `graphql:"createSchedule(input: $input)"`
}

type ScheduleMutateDeleteStruct struct {
	Schedule NewSchedule `graphql:"deleteSchedule(ID: $ID)"`
}

func (s *Schedule) Encode() (tf.M, error) {
	m, err := tf.Encode(s)
	if err != nil {
		return nil, err
	}

	m["team_id"] = s.Owner.ID

	return m, nil
}

// todo: encode tags
func (tag Tag) Encode() (tf.M, error) {
	return tf.Encode(tag)
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

func (rot Rotation) Encode() (tf.M, error) {
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

func (s *NewSchedule) Encode() (tf.M, error) {
	m, err := tf.Encode(s)
	if err != nil {
		return nil, err
	}

	m["team_id"] = s.TeamID

	tagsEncoded, terr := tf.EncodeSlice(s.Tags)
	if terr != nil {
		return nil, terr
	}
	m["tags"] = tagsEncoded

	rotationsEncoded, rerr := tf.EncodeSlice(s.Rotations)
	if rerr != nil {
		return nil, rerr
	}
	m["rotations"] = rotationsEncoded

	return m, nil
}

func (client *Client) GetScheduleById(ctx context.Context, teamID string, id string) (*Schedule, error) {
	url := fmt.Sprintf("%s/schedules/%s?owner_id=%s", client.BaseURLV3, id, teamID)

	return Request[any, Schedule](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) GetScheduleByName(ctx context.Context, teamID string, name string) (*Schedule, error) {
	schedules, err := client.ListSchedules(ctx, teamID)
	if err != nil {
		return nil, err
	}

	for _, s := range schedules {
		if s.Name == name {
			return s, nil
		}
	}

	return nil, fmt.Errorf("could not find a schedule with name `%s`", name)
}

func (client *Client) ListSchedules(ctx context.Context, teamID string) ([]*Schedule, error) {
	url := fmt.Sprintf("%s/schedules?owner_id=%s", client.BaseURLV3, teamID)

	return RequestSlice[any, Schedule](http.MethodGet, url, client, ctx, nil)
}

type CreateUpdateScheduleReq struct {
	Name        string `json:"name"`
	Color       string `json:"colour"`
	Description string `json:"description"`
	TeamID      string `json:"owner_id"`
}

func (client *Client) CreateSchedule(ctx context.Context, req *CreateUpdateScheduleReq) (*Schedule, error) {
	url := fmt.Sprintf("%s/schedules", client.BaseURLV3)

	return Request[CreateUpdateScheduleReq, Schedule](http.MethodPost, url, client, ctx, req)
}

func (client *Client) UpdateSchedule(ctx context.Context, id string, req *CreateUpdateScheduleReq) (*Schedule, error) {
	url := fmt.Sprintf("%s/schedules/%s", client.BaseURLV3, id)

	return Request[CreateUpdateScheduleReq, Schedule](http.MethodPut, url, client, ctx, req)
}

func (client *Client) DeleteSchedule(ctx context.Context, id string) (*any, error) {
	url := fmt.Sprintf("%s/schedules/%s", client.BaseURLV3, id)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}

// ScheduleV2 APIs
func (client *Client) DeleteScheduleV2ByID(ctx context.Context, ID string) (*ScheduleMutateDeleteStruct, error) {
	var m ScheduleMutateDeleteStruct

	id, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		diag.Errorf("unable to convert schedule ID to string")
	}

	variables := map[string]interface{}{
		"ID": id,
	}

	return GraphQLRequest[ScheduleMutateDeleteStruct]("mutate", client, ctx, &m, variables)
}

func (client *Client) GetScheduleV2ById(ctx context.Context, ID string) (*ScheduleQueryStruct, error) {
	var m ScheduleQueryStruct

	id, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		diag.Errorf("unable to convert schedule ID to string")
	}

	variables := map[string]interface{}{
		"ID": id,
	}

	return GraphQLRequest[ScheduleQueryStruct]("query", client, ctx, &m, variables)
}

func (client *Client) CreateScheduleV2(ctx context.Context, payload NewSchedule) (*ScheduleMutateStruct, error) {
	var m ScheduleMutateStruct

	variables := map[string]interface{}{
		"input": payload,
	}

	return GraphQLRequest[ScheduleMutateStruct]("mutate", client, ctx, &m, variables)
}
