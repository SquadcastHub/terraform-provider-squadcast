package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type GER struct {
	ID          uint         `json:"id" tf:"id"`
	TeamID      string       `json:"owner_id" tf:"team_id"`
	Name        string       `json:"name" tf:"name"`
	Description string       `json:"description" tf:"description"`
	RoutingKey  string       `json:"routing_key,omitempty" tf:"routing_key"`
	EntityOwner *EntityOwner `json:"entity_owner"`
}

type GER_Ruleset struct {
	ID                   uint              `json:"id" tf:"id"`
	GER_ID               uint              `json:"global_event_rule_id" tf:"ger_id"`
	AlertSourceName      string            `json:"alert_source" tf:"-"`
	AlertSourceShortName string            `json:"alert_source_shortname" tf:"alert_source_shortname"`
	AlertSourceVersion   string            `json:"alert_source_version" tf:"alert_source_version"`
	CatchAllAction       map[string]string `json:"catch_all_action" tf:"catch_all_action"`
	Ordering             []uint            `json:"ordering,omitempty" tf:"-"`
}

type GER_Ruleset_Rules struct {
	ID          uint              `json:"id" tf:"id"`
	GER_ID      uint              `json:"global_event_rule_id" tf:"ger_id"`
	Description string            `json:"description,omitempty" tf:"description"`
	Expression  string            `json:"expression,omitempty" tf:"expression"`
	Action      map[string]string `json:"action" tf:"action"`
}

type GERAlertSource struct {
	Name    string `tf:"name"`
	Version string `tf:"version"`
}

type GERReorderRulesetRulesReq struct {
	Ordering []uint `json:"ordering"`
}
type GERReorderRulesetRules struct {
	ID       uint   `json:"id,omitempty" tf:"id"`
	GER_ID   uint   `json:"global_event_rule_id" tf:"ger_id"`
	Ordering []uint `json:"ordering"`
}

func (ger *GER) Encode() (map[string]interface{}, error) {
	m, err := tf.Encode(ger)
	if err != nil {
		return nil, err
	}

	if ger.EntityOwner != nil {
		m["entity_owner"] = tf.List(tf.M{
			"id":   ger.EntityOwner.ID,
			"type": ger.EntityOwner.Type,
		})
	}

	return m, nil
}

func (ger *GER_Ruleset) Encode() (map[string]interface{}, error) {
	m, err := tf.Encode(ger)
	if err != nil {
		return nil, err
	}

	gerID := strconv.FormatUint(uint64(ger.GER_ID), 10)
	m["ger_id"] = gerID

	catchAllAction, err := tf.Encode(ger.CatchAllAction)
	if err != nil {
		return nil, err
	}
	m["catch_all_action"] = catchAllAction

	return m, nil
}

func (ger *GER_Ruleset_Rules) Encode() (map[string]interface{}, error) {
	m, err := tf.Encode(ger)
	if err != nil {
		return nil, err
	}

	gerID := strconv.FormatUint(uint64(ger.GER_ID), 10)
	m["ger_id"] = gerID

	action, err := tf.Encode(ger.Action)
	if err != nil {
		return nil, err
	}
	m["action"] = action

	return m, nil
}

func (ger *GERReorderRulesetRules) Encode() (map[string]interface{}, error) {
	m, err := tf.Encode(ger)
	if err != nil {
		return nil, err
	}

	gerID := strconv.FormatUint(uint64(ger.GER_ID), 10)
	m["ger_id"] = gerID

	ordering := make([]string, len(ger.Ordering))
	for i, v := range ger.Ordering {
		ordering[i] = strconv.FormatUint(uint64(v), 10)
	}
	m["ordering"] = ordering

	return m, nil
}

func (client *Client) CreateGER(ctx context.Context, req *GER) (*GER, error) {
	url := fmt.Sprintf("%s/global-event-rules", client.BaseURLV3)
	data, err := Request[GER, GER](http.MethodPost, url, client, ctx, req)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (client *Client) GetGERById(ctx context.Context, ID string) (*GER, error) {
	url := fmt.Sprintf("%s/global-event-rules/%s", client.BaseURLV3, ID)
	data, err := Request[any, GER](http.MethodGet, url, client, ctx, nil)
	if err != nil {
		return nil, errors.New("GER with given ID not found")
	}
	return data, nil
}

func (client *Client) UpdateGER(ctx context.Context, gerID string, req *GER) (*GER, error) {
	url := fmt.Sprintf("%s/global-event-rules/%s", client.BaseURLV3, gerID)
	return Request[GER, GER](http.MethodPatch, url, client, ctx, req)
}

func (client *Client) DeleteGER(ctx context.Context, ID string) (*any, error) {
	url := fmt.Sprintf("%s/global-event-rules/%s", client.BaseURLV3, ID)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}

// GER Ruleset APIs
func (client *Client) CreateGERRuleset(ctx context.Context, gerID string, req *GER_Ruleset) (*GER_Ruleset, error) {
	url := fmt.Sprintf("%s/global-event-rules/%s/rulesets", client.BaseURLV3, gerID)
	data, err := Request[GER_Ruleset, GER_Ruleset](http.MethodPost, url, client, ctx, req)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (client *Client) GetGERRulesetByAlertSource(ctx context.Context, gerID string, alertSource GERAlertSource) (*GER_Ruleset, error) {
	url := fmt.Sprintf("%s/global-event-rules/%s/rulesets/%s/%s", client.BaseURLV3, gerID, alertSource.Version, alertSource.Name)

	data, err := Request[any, GER_Ruleset](http.MethodGet, url, client, ctx, nil)
	if err != nil {
		return nil, errors.New(alertSource.Name + " GER Ruleset with given ID not found")
	}
	return data, nil
}

func (client *Client) UpdateGERRuleset(ctx context.Context, gerID string, alertSource GERAlertSource, req *GER_Ruleset) (*GER_Ruleset, error) {
	url := fmt.Sprintf("%s/global-event-rules/%s/rulesets/%s/%s", client.BaseURLV3, gerID, alertSource.Version, alertSource.Name)
	return Request[GER_Ruleset, GER_Ruleset](http.MethodPatch, url, client, ctx, req)
}

func (client *Client) DeleteGERRuleset(ctx context.Context, gerID string, alertSource GERAlertSource) (*any, error) {
	url := fmt.Sprintf("%s/global-event-rules/%s/rulesets/%s/%s", client.BaseURLV3, gerID, alertSource.Version, alertSource.Name)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}

// GER Ruleset Rules APIs
func (client *Client) CreateGERRulesetRules(ctx context.Context, gerID string, alertSource GERAlertSource, req *GER_Ruleset_Rules) (*GER_Ruleset_Rules, error) {
	url := fmt.Sprintf("%s/global-event-rules/%s/rulesets/%s/%s/rules", client.BaseURLV3, gerID, alertSource.Version, alertSource.Name)
	data, err := Request[GER_Ruleset_Rules, GER_Ruleset_Rules](http.MethodPost, url, client, ctx, req)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (client *Client) GetGERRulesetRulesById(ctx context.Context, gerID string, ruleID string, alertSource GERAlertSource) (*GER_Ruleset_Rules, error) {
	url := fmt.Sprintf("%s/global-event-rules/%s/rulesets/%s/%s/rules/%s", client.BaseURLV3, gerID, alertSource.Version, alertSource.Name, ruleID)
	data, err := Request[any, GER_Ruleset_Rules](http.MethodGet, url, client, ctx, nil)
	if err != nil {
		return nil, errors.New("GER Ruleset Rule with given ID not found")
	}
	return data, nil
}

func (client *Client) UpdateGERRulesetRules(ctx context.Context, gerID string, ruleID string, alertSource GERAlertSource, req *GER_Ruleset_Rules) (*GER_Ruleset_Rules, error) {
	url := fmt.Sprintf("%s/global-event-rules/%s/rulesets/%s/%s/rules/%s", client.BaseURLV3, gerID, alertSource.Version, alertSource.Name, ruleID)
	return Request[GER_Ruleset_Rules, GER_Ruleset_Rules](http.MethodPatch, url, client, ctx, req)
}

func (client *Client) DeleteGERRulesetRules(ctx context.Context, gerID string, ruleID string, alertSource GERAlertSource) (*any, error) {
	url := fmt.Sprintf("%s/global-event-rules/%s/rulesets/%s/%s/rules/%s", client.BaseURLV3, gerID, alertSource.Version, alertSource.Name, ruleID)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}

func (client *Client) UpdateGERRulesetRulesOrdering(ctx context.Context, gerID string, alertSource GERAlertSource, req *GERReorderRulesetRulesReq) (*GERReorderRulesetRules, error) {
	url := fmt.Sprintf("%s/global-event-rules/%s/rulesets/%s/%s/priority", client.BaseURLV3, gerID, alertSource.Version, alertSource.Name)
	return Request[GERReorderRulesetRulesReq, GERReorderRulesetRules](http.MethodPatch, url, client, ctx, req)
}
