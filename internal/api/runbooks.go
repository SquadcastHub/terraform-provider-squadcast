package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type RunbookStep struct {
	Content string `json:"content" tf:"content"`
}

func (rs *RunbookStep) Encode() (tf.M, error) {
	return tf.Encode(rs)
}

type Runbook struct {
	ID    string         `json:"id" tf:"id"`
	Name  string         `json:"name" tf:"name"`
	Steps []*RunbookStep `json:"steps" tf:"-"`
	Owner OwnerRef       `json:"owner" tf:"-"`
}

func (r *Runbook) Encode() (tf.M, error) {
	m, err := tf.Encode(r)
	if err != nil {
		return nil, err
	}

	steps, err := tf.EncodeSlice(r.Steps)
	if err != nil {
		return nil, err
	}
	m["steps"] = steps

	m["team_id"] = r.Owner.ID

	return m, nil
}

func (client *Client) GetRunbookById(ctx context.Context, teamID string, id string) (*Runbook, error) {
	url := fmt.Sprintf("%s/runbooks/%s?owner_id=%s", client.BaseURLV3, id, teamID)

	return Request[any, Runbook](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) GetRunbookByName(ctx context.Context, teamID string, name string) (*Runbook, error) {
	runbooks, err := client.ListRunbooks(ctx, teamID)
	if err != nil {
		return nil, err
	}

	for _, s := range runbooks {
		if s.Name == name {
			return s, nil
		}
	}

	return nil, fmt.Errorf("could not find a runbook with name `%s`", name)
}

func (client *Client) ListRunbooks(ctx context.Context, teamID string) ([]*Runbook, error) {
	url := fmt.Sprintf("%s/runbooks?owner_id=%s", client.BaseURLV3, teamID)

	return RequestSlice[any, Runbook](http.MethodGet, url, client, ctx, nil)
}

type CreateUpdateRunbookReq struct {
	Name   string         `json:"name"`
	TeamID string         `json:"owner_id"`
	Steps  []*RunbookStep `json:"steps"`
}

func (client *Client) CreateRunbook(ctx context.Context, req *CreateUpdateRunbookReq) (*Runbook, error) {
	url := fmt.Sprintf("%s/runbooks", client.BaseURLV3)

	return Request[CreateUpdateRunbookReq, Runbook](http.MethodPost, url, client, ctx, req)
}

func (client *Client) UpdateRunbook(ctx context.Context, id string, req *CreateUpdateRunbookReq) (*Runbook, error) {
	url := fmt.Sprintf("%s/runbooks/%s", client.BaseURLV3, id)

	return Request[CreateUpdateRunbookReq, Runbook](http.MethodPut, url, client, ctx, req)
}

func (client *Client) DeleteRunbook(ctx context.Context, id string) (*any, error) {
	url := fmt.Sprintf("%s/runbooks/%s", client.BaseURLV3, id)

	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
