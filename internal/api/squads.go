package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type OwnerRef struct {
	ID   string `json:"id" tf:"id"`
	Type string `json:"type" tf:"type"`
}

type Squad struct {
	ID        string   `json:"id" tf:"id"`
	Name      string   `json:"name" tf:"name"`
	Slug      string   `json:"slug" tf:"-"`
	Owner     OwnerRef `json:"owner" tf:"-"`
	MemberIDs []string `json:"members" tf:"member_ids"`
}

func (s *Squad) Encode() (tf.M, error) {
	m, err := tf.Encode(s)
	if err != nil {
		return nil, err
	}

	m["team_id"] = s.Owner.ID

	return m, nil
}

func (client *Client) GetSquadById(ctx context.Context, teamID string, id string) (*Squad, error) {
	url := fmt.Sprintf("%s/squads/%s?owner_id=%s", client.BaseURLV3, id, teamID)

	return Request[any, Squad](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) GetSquadByName(ctx context.Context, teamID string, name string) (*Squad, error) {
	url := fmt.Sprintf("%s/squads/by-name?name=%s&owner_id=%s", client.BaseURLV3, url.QueryEscape(name), teamID)

	return Request[any, Squad](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) ListSquads(ctx context.Context, teamID string) ([]*Squad, error) {
	url := fmt.Sprintf("%s/squads?owner_id=%s", client.BaseURLV3, teamID)

	return RequestSlice[any, Squad](http.MethodGet, url, client, ctx, nil)
}

type CreateSquadReq struct {
	Name      string   `json:"name"`
	TeamID    string   `json:"owner_id"`
	MemberIDs []string `json:"members"`
}

type UpdateSquadReq struct {
	Name      string   `json:"name"`
	MemberIDs []string `json:"members"`
}

func (client *Client) CreateSquad(ctx context.Context, req *CreateSquadReq) (*Squad, error) {
	url := fmt.Sprintf("%s/squads", client.BaseURLV3)

	return Request[CreateSquadReq, Squad](http.MethodPost, url, client, ctx, req)
}

func (client *Client) UpdateSquad(ctx context.Context, id string, req *UpdateSquadReq) (*Squad, error) {
	url := fmt.Sprintf("%s/squads/%s", client.BaseURLV3, id)

	return Request[UpdateSquadReq, Squad](http.MethodPut, url, client, ctx, req)
}

func (client *Client) DeleteSquad(ctx context.Context, id string) (*any, error) {
	url := fmt.Sprintf("%s/squads/%s", client.BaseURLV3, id)

	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
