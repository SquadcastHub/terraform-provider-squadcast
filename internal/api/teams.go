package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type Team struct {
	ID          string            `json:"id" tf:"id"`
	Name        string            `json:"name" tf:"name"`
	Description string            `json:"description" tf:"description"`
	Default     bool              `json:"default" tf:"default"`
	Members     []*DataTeamMember `json:"members" tf:"-"`
	Roles       []*TeamRole       `json:"roles" tf:"-"`
}

func (t *Team) Encode() (tf.M, error) {
	m, err := tf.Encode(t)
	if err != nil {
		return nil, err
	}

	members, err := tf.EncodeSlice(t.Members)
	if err != nil {
		return nil, err
	}
	m["members"] = members

	roles, err := tf.EncodeSlice(t.Roles)
	if err != nil {
		return nil, err
	}
	m["roles"] = roles

	return m, nil
}

type DataTeamMember struct {
	UserID  string   `json:"user_id" tf:"user_id"`
	RoleIDs []string `json:"role_ids" tf:"role_ids"`
}

func (tm *DataTeamMember) Encode() (tf.M, error) {
	return tf.Encode(tm)
}

type TeamMember struct {
	ID      string   `tf:"id"`
	UserID  string   `json:"user_id" tf:"user_id"`
	RoleIDs []string `json:"role_ids" tf:"role_ids"`
}

func (tm *TeamMember) Encode() (tf.M, error) {
	tm.ID = tm.UserID
	return tf.Encode(tm)
}

type TeamRole struct {
	ID        string                 `json:"id" tf:"id"`
	Name      string                 `json:"name" tf:"name"`
	Slug      string                 `json:"slug" tf:"-"`
	Default   bool                   `json:"default" tf:"default"`
	Abilities RBACEntityAbilitiesMap `json:"abilities" tf:"-"`
}

func (tr *TeamRole) Encode() (tf.M, error) {
	m, err := tf.Encode(tr)
	if err != nil {
		return nil, err
	}

	abilities := make([]string, 0, 100)
	for _, kv := range tr.Abilities {
		for k := range kv {
			abilities = append(abilities, k)
		}
	}

	sort.Strings(abilities)
	m["abilities"] = abilities

	return m, nil
}

type RBACAbilityMap map[string]bool
type RBACEntityAbilitiesMap map[string]RBACAbilityMap

func (client *Client) GetTeamByName(ctx context.Context, name string) (*Team, error) {
	url := fmt.Sprintf("%s/teams/by-name?name=%s", client.BaseURLV3, url.QueryEscape(name))

	return Request[any, Team](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) GetTeamById(ctx context.Context, id string) (*Team, error) {
	url := fmt.Sprintf("%s/teams/%s", client.BaseURLV3, id)

	return Request[any, Team](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) DeleteTeam(ctx context.Context, id string) (*any, error) {
	url := fmt.Sprintf("%s/teams/%s", client.BaseURLV3, id)

	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
