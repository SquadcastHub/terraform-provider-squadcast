package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type DeduplicationRuleCondition struct {
	LHS string `json:"lhs" tf:"lhs"`
	Op  string `json:"op" tf:"op"`
	RHS string `json:"rhs" tf:"rhs"`
}

func (c *DeduplicationRuleCondition) Encode() (tf.M, error) {
	return tf.Encode(c)
}

type DeduplicationRule struct {
	ID                      string                        `json:"rule_id,omitempty" tf:"-"`
	IsBasic                 bool                          `json:"is_basic" tf:"is_basic"`
	Description             string                        `json:"description" tf:"description"`
	Expression              string                        `json:"expression" tf:"expression"`
	DependencyDeduplication bool                          `json:"dependency_deduplication" tf:"dependency_deduplication"`
	TimeUnit                string                        `json:"time_unit" tf:"time_unit"`
	TimeWindow              int                           `json:"time_window" tf:"time_window"`
	BasicExpression         []*DeduplicationRuleCondition `json:"basic_expression" tf:"basic_expressions"`
}

func (r *DeduplicationRule) Encode() (tf.M, error) {
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

type DeduplicationRules struct {
	ID        string               `json:"id" tf:"id"`
	ServiceID string               `json:"service_id" tf:"service_id"`
	Rules     []*DeduplicationRule `json:"rules" tf:"-"`
}

func (s *DeduplicationRules) Encode() (tf.M, error) {
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

func (client *Client) GetDeduplicationRules(ctx context.Context, serviceID, teamID string) (*DeduplicationRules, error) {
	url := fmt.Sprintf("%s/services/%s/deduplication-rules?owner_id=%s", client.BaseURLV3, serviceID, teamID)

	return Request[any, DeduplicationRules](http.MethodGet, url, client, ctx, nil)
}

type UpdateDeduplicationRulesReq struct {
	Rules []DeduplicationRule `json:"rules"`
}

func (client *Client) UpdateDeduplicationRules(ctx context.Context, serviceID, teamID string, req *UpdateDeduplicationRulesReq) (*DeduplicationRules, error) {
	url := fmt.Sprintf("%s/services/%s/deduplication-rules", client.BaseURLV3, serviceID)
	return Request[UpdateDeduplicationRulesReq, DeduplicationRules](http.MethodPost, url, client, ctx, req)
}

// deduplication rules v2

type CreateDeduplicationRule struct {
	Rule DeduplicationRule `json:"rule"`
}
type DeduplicationRuleV2 struct {
	ID        string             `json:"rule_id" tf:"id"`
	ServiceID string             `json:"service_id" tf:"service_id"`
	Rule      *DeduplicationRule `json:"rule" tf:"-"`
}

func (s *DeduplicationRuleV2) Encode() (tf.M, error) {
	m, err := tf.Encode(s.Rule)
	if err != nil {
		return nil, err
	}

	basicExpressions, err := tf.EncodeSlice(s.Rule.BasicExpression)
	if err != nil {
		return nil, err
	}
	m["basic_expressions"] = basicExpressions

	return m, nil
}

func (client *Client) CreateDeduplicationRulesV2(ctx context.Context, serviceID string, req *CreateDeduplicationRule) (*CreateDeduplicationRule, error) {
	url := fmt.Sprintf("%s/services/%s/deduplication-rules/new", client.BaseURLV3, serviceID)
	return Request[CreateDeduplicationRule, CreateDeduplicationRule](http.MethodPost, url, client, ctx, req)
}

func (client *Client) UpdateDeduplicationRuleByID(ctx context.Context, serviceID, ruleID string, req *CreateDeduplicationRule) (*CreateDeduplicationRule, error) {
	url := fmt.Sprintf("%s/services/%s/deduplication-rules/%s", client.BaseURLV3, serviceID, ruleID)
	return Request[CreateDeduplicationRule, CreateDeduplicationRule](http.MethodPut, url, client, ctx, req)
}

func (client *Client) GetDeduplicationRuleByID(ctx context.Context, serviceID, ruleID string) (*DeduplicationRuleV2, error) {
	url := fmt.Sprintf("%s/services/%s/deduplication-rules/%s", client.BaseURLV3, serviceID, ruleID)

	return Request[any, DeduplicationRuleV2](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) DeleteDeduplicationRuleByID(ctx context.Context, serviceID, ruleID string) (any, error) {
	url := fmt.Sprintf("%s/services/%s/deduplication-rules/%s", client.BaseURLV3, serviceID, ruleID)

	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
