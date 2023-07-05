package api

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type NewSchedule struct {
	ID          int    `graphql:"ID" json:"ID,omitempty" tf:"id"`
	Name        string `graphql:"name" json:"name" tf:"name"`
	Description string `graphql:"description" json:"description,omitempty" tf:"description"`
	TimeZone    string `graphql:"timeZone" json:"timeZone" tf:"timezone"`
	TeamID      string `graphql:"teamID" json:"teamID" tf:"team_id"`
	Owner       *Owner `graphql:"owner" json:"owner" tf:"-"`
	Tags        []*Tag `graphql:"tags" json:"tags,omitempty" tf:"tags"`
}

type UpdateSchedule struct {
	Name        string `graphql:"name" json:"name" tf:"name"`
	Description string `graphql:"description" json:"description,omitempty" tf:"description"`
	TimeZone    string `graphql:"timeZone" json:"timeZone" tf:"timezone"`
	Owner       *Owner `graphql:"owner" json:"owner" tf:"-"`
	Tags        []*Tag `graphql:"tags" json:"tags,omitempty" tf:"tags"`
}

type Owner struct {
	ID   string `graphql:"ID" json:"ID" tf:"id"`
	Type string `graphql:"type" json:"type" tf:"type"`
}

type Tag struct {
	Key   string `graphql:"key" json:"key" tf:"key"`
	Value string `graphql:"value" json:"value" tf:"value"`
	Color string `graphql:"color" json:"color,omitempty" tf:"color"`
}

// GraphQL query structs
type ScheduleQueryStruct struct {
	NewSchedule `graphql:"schedule(ID: $ID)"`
}

type ScheduleByNameQueryStruct struct {
	NewSchedule []*NewSchedule `graphql:"schedules(filters:  { scheduleName: $scheduleName, teamID: $teamID })"`
}

type CreateScheduleMutateStruct struct {
	NewSchedule `graphql:"createSchedule(input: $input)"`
}

type UpdateScheduleMutateStruct struct {
	UpdateSchedule `graphql:"updateSchedule(ID: $ID, input: $input)"`
}

type DeleteScheduleResponse struct {
	ID   int    `graphql:"ID"`
	Name string `graphql:"name"`
}
type ScheduleMutateDeleteStruct struct {
	Schedule DeleteScheduleResponse `graphql:"deleteSchedule(ID: $ID)"`
}

func (tag Tag) Encode() (tf.M, error) {
	return tf.Encode(tag)
}

func (s *NewSchedule) Encode() (tf.M, error) {
	m, err := tf.Encode(s)
	if err != nil {
		return nil, err
	}
	m["id"] = strconv.Itoa(s.ID)

	tagsEncoded, terr := tf.EncodeSlice(s.Tags)
	if terr != nil {
		return nil, terr
	}
	m["tags"] = tagsEncoded

	m["entity_owner"] = tf.List(tf.M{
		"id":   s.Owner.ID,
		"type": s.Owner.Type,
	})

	return m, nil
}

type CreateUpdateScheduleReq struct {
	Name        string `json:"name"`
	Color       string `json:"colour"`
	Description string `json:"description"`
	TeamID      string `json:"owner_id"`
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

func (client *Client) CreateScheduleV2(ctx context.Context, payload NewSchedule) (*CreateScheduleMutateStruct, error) {
	var m CreateScheduleMutateStruct

	variables := map[string]interface{}{
		"input": payload,
	}

	return GraphQLRequest[CreateScheduleMutateStruct]("mutate", client, ctx, &m, variables)
}

func (client *Client) UpdateScheduleV2(ctx context.Context, ID int, payload UpdateSchedule) (*UpdateScheduleMutateStruct, error) {
	var m UpdateScheduleMutateStruct

	variables := map[string]interface{}{
		"ID":    ID,
		"input": payload,
	}

	return GraphQLRequest[UpdateScheduleMutateStruct]("mutate", client, ctx, &m, variables)
}

func (client *Client) GetScheduleV2ByName(ctx context.Context, teamID string, scheduleName string) (*ScheduleByNameQueryStruct, error) {
	var m ScheduleByNameQueryStruct

	variables := map[string]interface{}{
		"scheduleName": scheduleName,
		"teamID":       teamID,
	}

	return GraphQLRequest[ScheduleByNameQueryStruct]("query", client, ctx, &m, variables)
}
