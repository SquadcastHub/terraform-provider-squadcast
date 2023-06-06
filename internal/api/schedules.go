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
	ID          int    `graphql:"ID" json:"id,omitempty" tf:"id"`
	Name        string `graphql:"name" json:"name" tf:"name"`
	Description string `graphql:"description" json:"description" tf:"description"`
	TimeZone    string `graphql:"timeZone" json:"timeZone" tf:"timezone"`
	TeamID      string `graphql:"teamID" json:"teamID" tf:"team_id"`
	// Tags        Tags     `graphql:"tags" json:"tags"`
	Owner *Owner `graphql:"owner" json:"owner" tf:"-"`
}

type Owner struct {
	ID   string `graphql:"ID" json:"ID" tf:"id"`
	Type string `graphql:"type" json:"type" tf:"type"`
}

type Tags struct {
	ID   string `graphql:"ID" json:"id" tf:"id"`
	Type string `graphql:"type" json:"type" tf:"type"`
}

// GraphQL query structs
type ScheduleQueryStruct struct {
	NewSchedule `graphql:"schedule(ID: $ID)"`
}

type ScheduleMutateStruct struct {
	NewSchedule `graphql:"createSchedule(input: $input)"`
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
func (s *NewSchedule) Encode() (tf.M, error) {
	m, err := tf.Encode(s)
	if err != nil {
		return nil, err
	}

	m["team_id"] = s.Owner.ID

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
