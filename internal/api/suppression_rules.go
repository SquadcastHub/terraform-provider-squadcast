package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type SuppressionRuleCondition struct {
	LHS string `json:"lhs" tf:"lhs"`
	Op  string `json:"op" tf:"op"`
	RHS string `json:"rhs" tf:"rhs"`
}

func (c *SuppressionRuleCondition) Encode() (tf.M, error) {
	return tf.Encode(c)
}

type SuppressionRule struct {
	IsBasic         bool                        `json:"is_basic" tf:"is_basic"`
	Description     string                      `json:"description" tf:"description"`
	Expression      string                      `json:"expression" tf:"expression"`
	BasicExpression []*SuppressionRuleCondition `json:"basic_expression" tf:"basic_expressions"`
}

func (r *SuppressionRule) Encode() (tf.M, error) {
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

type SuppressionRules struct {
	ID        string             `json:"id" tf:"id"`
	ServiceID string             `json:"service_id" tf:"service_id"`
	Rules     []*SuppressionRule `json:"rules" tf:"-"`
}

func (s *SuppressionRules) Encode() (tf.M, error) {
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

func (client *Client) GetSuppressionRules(ctx context.Context, serviceID, teamID string) (*SuppressionRules, error) {
	url := fmt.Sprintf("%s/services/%s/suppression-rules", client.BaseURLV3, serviceID)

	return Request[any, SuppressionRules](http.MethodGet, url, client, ctx, nil)
}

type UpdateSuppressionRulesReq struct {
	Rules []SuppressionRule `json:"rules"`
}

func (client *Client) UpdateSuppressionRules(ctx context.Context, serviceID, teamID string, req *UpdateSuppressionRulesReq) (*SuppressionRules, error) {
	url := fmt.Sprintf("%s/services/%s/suppression-rules", client.BaseURLV3, serviceID)
	return Request[UpdateSuppressionRulesReq, SuppressionRules](http.MethodPost, url, client, ctx, req)
}
