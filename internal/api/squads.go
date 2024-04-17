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
	ID      string   `json:"id" tf:"id"`
	Name    string   `json:"name" tf:"name"`
	TeamID  string   `json:"team_id" tf:"team_id"`
	Members []Member `json:"members" tf:"members"`
}

func (m Member) Encode() (tf.M, error) {
	return tf.Encode(m)
}

func (s *Squad) Encode() (tf.M, error) {
	m, err := tf.Encode(s)
	if err != nil {
		return nil, err
	}

	members, err := tf.EncodeSlice(s.Members)
	if err != nil {
		return nil, err
	}

	m["members"] = members

	return m, nil
}

func (client *Client) GetSquadById(ctx context.Context, id string) (*Squad, error) {
	url := fmt.Sprintf("%s/squads/%s", client.BaseURLV4, id)

	return Request[any, Squad](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) GetSquadByName(ctx context.Context, teamID string, name string) (*Squad, error) {
	url := fmt.Sprintf("%s/squads/by-name?name=%s&owner_id=%s", client.BaseURLV3, url.QueryEscape(name), teamID)

	return Request[any, Squad](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) ListSquads(ctx context.Context, teamID string) ([]*Squad, error) {
	url := fmt.Sprintf("%s/squads?team_id=%s", client.BaseURLV4, teamID)

	return RequestSlice[any, Squad](http.MethodGet, url, client, ctx, nil)
}

type CreateSquadReq struct {
	Name   string `json:"name"`
	TeamID string `json:"owner_id"`
	// MemberIDs []string `json:"members"`
	Members []Member `json:"members"`
}

type Member struct {
	UserID string `json:"user_id" tf:"user_id"`
	Role   string `json:"role,omitempty" tf:"role"`
}

type UpdateSquadReq struct {
	Name string `json:"name"`
}

func (client *Client) CreateSquad(ctx context.Context, req *CreateSquadReq) (*Squad, error) {
	url := fmt.Sprintf("%s/squads", client.BaseURLV4)

	return Request[CreateSquadReq, Squad](http.MethodPost, url, client, ctx, req)
}

func (client *Client) UpdateSquad(ctx context.Context, id string, req *UpdateSquadReq) (*Squad, error) {
	url := fmt.Sprintf("%s/squads/%s/name", client.BaseURLV4, id)

	return Request[UpdateSquadReq, Squad](http.MethodPut, url, client, ctx, req)
}

type AddSquadMemberReq struct {
	Members []Member `json:"members"`
}

func (client *Client) AddSquadMembers(ctx context.Context, id string, req *AddSquadMemberReq) (*Squad, error) {
	url := fmt.Sprintf("%s/squads/%s/members", client.BaseURLV4, id)

	return Request[AddSquadMemberReq, Squad](http.MethodPost, url, client, ctx, req)
}

func (client *Client) RemoveSquadMember(ctx context.Context, id, memberID string) (any, error) {
	url := fmt.Sprintf("%s/squads/%s/members/%s", client.BaseURLV4, id, memberID)

	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}

func (client *Client) DeleteSquad(ctx context.Context, id string) (*any, error) {
	url := fmt.Sprintf("%s/squads/%s", client.BaseURLV3, id)

	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
