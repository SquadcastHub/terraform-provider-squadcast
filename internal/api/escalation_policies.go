package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type EscalationPolicyTarget struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`
	PID  int    `json:"pid,omitempty"`
}

func (t *EscalationPolicyTarget) Encode() (tf.M, error) {
	var ID string
	if t.Type == "schedulev2" {
		ID = fmt.Sprintf("%d", t.PID)
	} else {
		ID = t.ID
	}
	return tf.M{
		"id":   ID,
		"type": t.Type,
	}, nil
}

type EscalationPolicyRule struct {
	EscalateAfterMinutes     int                       `json:"escalationTime"`
	Via                      []string                  `json:"via"`
	Targets                  []*EscalationPolicyTarget `json:"entities"`
	RoundrobinEnabled        bool                      `json:"roundrobin_enabled"`
	EscalateWithinRoundrobin bool                      `json:"escalate_within_roundrobin"`
	RepeatTimes              int                       `json:"repetition"`
	RepeatAfterMinutes       int                       `json:"repeat_after"`
}

func (r *EscalationPolicyRule) Encode() (tf.M, error) {
	m := tf.M{
		"delay_minutes": r.EscalateAfterMinutes,
	}

	if len(r.Via) == 0 {
		m["notification_channels"] = []string{}
	} else {
		m["notification_channels"] = r.Via
	}

	if !r.RoundrobinEnabled || !r.EscalateWithinRoundrobin {
		if r.RepeatTimes != 0 || r.RepeatAfterMinutes != 0 {
			m["repeat"] = tf.List(tf.M{
				"times":         r.RepeatTimes,
				"delay_minutes": r.RepeatAfterMinutes,
			})
		}
	}

	if r.RoundrobinEnabled {
		rr := tf.M{
			"enabled": true,
		}

		if r.EscalateWithinRoundrobin {
			rr["rotation"] = tf.List(tf.M{
				"enabled":       true,
				"delay_minutes": r.RepeatAfterMinutes,
			})
		}

		m["round_robin"] = tf.List(rr)
	}

	targets, err := tf.EncodeSlice(r.Targets)
	if err != nil {
		return nil, err
	}
	m["targets"] = targets

	return m, nil
}

type EscalationPolicy struct {
	ID                 string                  `json:"id"`
	Name               string                  `json:"name"`
	Description        string                  `json:"description"`
	RepeatTimes        int                     `json:"repetition"`
	RepeatAfterMinutes int                     `json:"repeat_after"`
	Rules              []*EscalationPolicyRule `json:"rules"`
	Slug               string                  `json:"slug"`
	Owner              OwnerRef                `json:"owner"`
	EntityOwner        *EntityOwner            `json:"entity_owner"`
}

func (ep *EscalationPolicy) Encode() (tf.M, error) {
	m := tf.M{
		"id":          ep.ID,
		"name":        ep.Name,
		"description": ep.Description,
		"team_id":     ep.Owner.ID,
	}

	if ep.RepeatTimes != 0 || ep.RepeatAfterMinutes != 0 {
		m["repeat"] = tf.List(tf.M{
			"times":         ep.RepeatTimes,
			"delay_minutes": ep.RepeatAfterMinutes,
		})
	}

	rules, err := tf.EncodeSlice(ep.Rules)
	if err != nil {
		return nil, err
	}

	m["rules"] = rules

	if ep.EntityOwner != nil {
		m["entity_owner"] = tf.List(tf.M{
			"id":   ep.EntityOwner.ID,
			"type": ep.EntityOwner.Type,
		})
	}
	return m, nil
}

func (client *Client) GetEscalationPolicyById(ctx context.Context, teamID string, id string) (*EscalationPolicy, error) {
	url := fmt.Sprintf("%s/escalation-policies/%s?owner_id=%s", client.BaseURLV3, id, teamID)

	return Request[any, EscalationPolicy](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) GetEscalationPolicyByName(ctx context.Context, teamID string, name string) (*EscalationPolicy, error) {
	escalationPolicies, err := client.ListEscalationPolicies(ctx, teamID)
	if err != nil {
		return nil, err
	}

	for _, s := range escalationPolicies {
		if s.Name == name {
			return s, nil
		}
	}

	return nil, fmt.Errorf("could not find an escalation policy with name `%s`", name)
}

func (client *Client) ListEscalationPolicies(ctx context.Context, teamID string) ([]*EscalationPolicy, error) {
	url := fmt.Sprintf("%s/escalation-policies?owner_id=%s", client.BaseURLV3, teamID)

	return RequestSlice[any, EscalationPolicy](http.MethodGet, url, client, ctx, nil)
}

type CreateUpdateEscalationPolicyReq struct {
	TeamID             string                 `json:"owner_id"`
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	RepeatTimes        int                    `json:"repetition"`
	RepeatAfterMinutes int                    `json:"repeat_after"`
	Rules              []EscalationPolicyRule `json:"rules"`
	IsUsingNewFields   bool                   `json:"is_using_new_fields"`
	EntityOwner        *EntityOwner           `json:"entity_owner"`
}

func (client *Client) CreateEscalationPolicy(ctx context.Context, req *CreateUpdateEscalationPolicyReq) (*EscalationPolicy, error) {
	url := fmt.Sprintf("%s/escalation-policies", client.BaseURLV3)

	return Request[CreateUpdateEscalationPolicyReq, EscalationPolicy](http.MethodPost, url, client, ctx, req)
}

func (client *Client) UpdateEscalationPolicy(ctx context.Context, id string, req *CreateUpdateEscalationPolicyReq) (*EscalationPolicy, error) {
	url := fmt.Sprintf("%s/escalation-policies/%s", client.BaseURLV3, id)

	return Request[CreateUpdateEscalationPolicyReq, EscalationPolicy](http.MethodPost, url, client, ctx, req)
}

func (client *Client) DeleteEscalationPolicy(ctx context.Context, id string) (*any, error) {
	url := fmt.Sprintf("%s/escalation-policies/%s", client.BaseURLV3, id)

	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
