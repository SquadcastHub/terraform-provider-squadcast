package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func (client *Client) ListTeamRoles(ctx context.Context, teamID string) ([]*TeamRole, error) {
	url := fmt.Sprintf("%s/teams/%s/roles", client.BaseURLV3, teamID)

	return RequestSlice[any, TeamRole](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) GetTeamRoleByID(ctx context.Context, teamID string, id string) (*TeamRole, error) {
	teamRoles, err := client.ListTeamRoles(ctx, teamID)
	if err != nil {
		return nil, err
	}
	for _, teamRole := range teamRoles {
		if teamRole.ID == id {
			return teamRole, nil
		}
	}
	return nil, fmt.Errorf("[404] could not find team role with the id: %s", id)
}

func (client *Client) GetTeamRoleByName(ctx context.Context, teamID string, name string) (*TeamRole, error) {
	teamRoles, err := client.ListTeamRoles(ctx, teamID)
	if err != nil {
		return nil, err
	}
	for _, teamRole := range teamRoles {
		if teamRole.Name == name {
			return teamRole, nil
		}
	}
	return nil, fmt.Errorf("GetTeamRoleByName: could not find team role with the name: %s", name)
}

type CreateTeamRoleReq struct {
	Name      string
	Abilities []string
}

type UpdateTeamRoleReq struct {
	Name      string
	Abilities []string
}

func decodeAbilities(ab []string) tf.M {
	abilities := tf.M{}
	for _, abilityStr := range ab {
		parts := strings.Split(abilityStr, "-")
		_, entityParts := parts[0], parts[1:]
		entity := strings.Join(entityParts, "_")
		entitymap, ok := abilities[entity]
		if !ok {
			abilities[entity] = tf.M{}
			entitymap = abilities[entity]
		}
		entitymap.(tf.M)[abilityStr] = true
	}

	return abilities
}

func (client *Client) CreateTeamRole(ctx context.Context, teamID string, req *CreateTeamRoleReq) (*TeamRole, error) {
	url := fmt.Sprintf("%s/teams/%s/roles", client.BaseURLV3, teamID)

	payload := tf.M{}
	payload["name"] = req.Name
	payload["abilities"] = decodeAbilities(req.Abilities)

	_, err := Request[tf.M, Team](http.MethodPost, url, client, ctx, &payload)
	if err != nil {
		return nil, err
	}

	return client.GetTeamRoleByName(ctx, teamID, req.Name)
}

func (client *Client) UpdateTeamRole(ctx context.Context, teamID string, id string, req *UpdateTeamRoleReq) (*TeamRole, error) {
	url := fmt.Sprintf("%s/teams/%s/roles/%s", client.BaseURLV3, teamID, id)

	payload := tf.M{}
	payload["name"] = req.Name
	payload["abilities"] = decodeAbilities(req.Abilities)

	_, err := Request[tf.M, Team](http.MethodPut, url, client, ctx, &payload)
	if err != nil {
		return nil, err
	}

	return client.GetTeamRoleByID(ctx, teamID, id)
}

func (client *Client) DeleteTeamRole(ctx context.Context, teamID string, id string) (*any, error) {
	url := fmt.Sprintf("%s/teams/%s/roles/%s", client.BaseURLV3, teamID, id)

	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
