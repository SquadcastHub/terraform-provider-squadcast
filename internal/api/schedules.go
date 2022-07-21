package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type Schedule struct {
	ID          string   `json:"id" tf:"id"`
	Name        string   `json:"name" tf:"name"`
	Slug        string   `json:"slug" tf:"-"`
	Colour      string   `json:"colour" tf:"color"`
	Description string   `json:"description" tf:"description"`
	Owner       OwnerRef `json:"owner" tf:"-"`
}

func (s *Schedule) Encode() (tf.M, error) {
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
