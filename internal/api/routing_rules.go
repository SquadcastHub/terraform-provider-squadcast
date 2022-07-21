package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type RoutingRuleCondition struct {
	LHS string `json:"lhs" tf:"lhs"`
	RHS string `json:"rhs" tf:"rhs"`
}

func (c *RoutingRuleCondition) Encode() (tf.M, error) {
	return tf.Encode(c)
}

type RouteTo struct {
	EntityID   string `json:"entity_id" tf:"route_to_id"`
	EntityType string `json:"entity_type" tf:"route_to_type"`
}

type RoutingRule struct {
	IsBasic         bool                    `json:"is_basic" tf:"is_basic"`
	Expression      string                  `json:"expression" tf:"expression"`
	BasicExpression []*RoutingRuleCondition `json:"basic_expression" tf:"basic_expressions"`
	RouteTo         RouteTo                 `json:"route_to" tf:"route_to,squash"`
}

func (r *RoutingRule) Encode() (tf.M, error) {
	m, err := tf.Encode(r)
	if err != nil {
		return nil, err
	}

	basicExpression, err := tf.EncodeSlice(r.BasicExpression)
	if err != nil {
		return nil, err
	}
	m["basic_expressions"] = basicExpression

	return m, nil
}

type RoutingRules struct {
	ID        string         `json:"id" tf:"id"`
	ServiceID string         `json:"service_id" tf:"service_id"`
	Rules     []*RoutingRule `json:"rules" tf:"-"`
}

func (s *RoutingRules) Encode() (tf.M, error) {
	m, err := tf.Encode(s)
	if err != nil {
		return nil, err
	}

	rules, err := tf.EncodeSlice(s.Rules)
	if err != nil {
		return nil, err
	}
	m["rules"] = rules

	return m, nil
}

func (client *Client) GetRoutingRules(ctx context.Context, serviceID, teamID string) (*RoutingRules, error) {
	url := fmt.Sprintf("%s/services/%s/routing-rules", client.BaseURLV3, serviceID)

	return Request[any, RoutingRules](http.MethodGet, url, client, ctx, nil)
}

type UpdateRoutingRulesReq struct {
	Rules []RoutingRule `json:"rules"`
}

func (client *Client) UpdateRoutingRules(ctx context.Context, serviceID, teamID string, req *UpdateRoutingRulesReq) (*RoutingRules, error) {
	url := fmt.Sprintf("%s/services/%s/routing-rules", client.BaseURLV3, serviceID)
	return Request[UpdateRoutingRulesReq, RoutingRules](http.MethodPost, url, client, ctx, req)
}
