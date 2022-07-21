package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type teamMetaRole struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Default bool   `json:"default"`
}

type TeamMeta struct {
	ID          string          `json:"id" tf:"id"`
	Name        string          `json:"name" tf:"name"`
	Description string          `json:"description" tf:"description"`
	Default     bool            `json:"default" tf:"default"`
	Roles       []*teamMetaRole `json:"roles" tf:"-"`
}

func (t *TeamMeta) Encode() (tf.M, error) {
	m, err := tf.Encode(t)
	if err != nil {
		return nil, err
	}

	defaultRoleNames := map[string]string{
		"Manage Team": "manage_team",
		"Admin":       "admin",
		"User":        "user",
		"Observer":    "observer",
	}

	roles := tf.M{}

	for _, role := range t.Roles {
		key := defaultRoleNames[role.Name]
		if key != "" {
			roles[key] = role.ID
		}
	}
	m["default_role_ids"] = roles

	return m, nil
}

func (client *Client) GetTeamMetaById(ctx context.Context, id string) (*TeamMeta, error) {
	team, err := client.GetTeamById(ctx, id)
	if err != nil {
		return nil, err
	}

	roles := make([]*teamMetaRole, len(team.Roles))

	for i, v := range team.Roles {
		roles[i] = &teamMetaRole{
			ID:      v.ID,
			Name:    v.Name,
			Default: v.Default,
		}
	}

	return &TeamMeta{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		Default:     team.Default,
		Roles:       roles,
	}, nil

}

type CreateTeamReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (client *Client) CreateTeam(ctx context.Context, req *CreateTeamReq) (*TeamMeta, error) {
	url := fmt.Sprintf("%s/teams", client.BaseURLV3)

	return Request[CreateTeamReq, TeamMeta](http.MethodPost, url, client, ctx, req)
}

type UpdateTeamMetaReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (client *Client) UpdateTeamMeta(ctx context.Context, id string, req *UpdateTeamMetaReq) (*TeamMeta, error) {
	url := fmt.Sprintf("%s/teams/%s/meta", client.BaseURLV3, id)

	return Request[UpdateTeamMetaReq, TeamMeta](http.MethodPatch, url, client, ctx, req)
}
