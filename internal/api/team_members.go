package api

import (
	"context"
	"fmt"
	"net/http"
)

func (client *Client) GetTeamMemberByID(ctx context.Context, teamID string, userID string) (*TeamMember, error) {
	url := fmt.Sprintf("%s/teams/%s/members/%s?owner_id=%s", client.BaseURLV3, teamID, userID, teamID)

	return Request[any, TeamMember](http.MethodGet, url, client, ctx, nil)
}

type CreateTeamMemberReq struct {
	UserID  string   `json:"user_id" tf:"user_id"`
	RoleIDs []string `json:"role_ids" tf:"role_ids"`
}

type UpdateTeamMemberReq struct {
	RoleIDs []string `json:"role_ids" tf:"role_ids"`
}

func (client *Client) CreateTeamMember(ctx context.Context, teamID string, req *CreateTeamMemberReq) (*TeamMember, error) {
	url := fmt.Sprintf("%s/teams/%s/members?owner_id=%s", client.BaseURLV3, teamID, teamID)

	return Request[CreateTeamMemberReq, TeamMember](http.MethodPost, url, client, ctx, req)
}

func (client *Client) UpdateTeamMember(ctx context.Context, teamID string, userID string, req *UpdateTeamMemberReq) (*TeamMember, error) {
	url := fmt.Sprintf("%s/teams/%s/members/%s?owner_id=%s", client.BaseURLV3, teamID, userID, teamID)

	return Request[UpdateTeamMemberReq, TeamMember](http.MethodPatch, url, client, ctx, req)
}

func (client *Client) DeleteTeamMember(ctx context.Context, teamID string, userID string) (*any, error) {
	url := fmt.Sprintf("%s/teams/%s/members/%s?owner_id=%s", client.BaseURLV3, teamID, userID, teamID)

	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
