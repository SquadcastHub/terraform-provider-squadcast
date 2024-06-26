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
	ID              string                      `json:"rule_id,omitempty" tf:"-"`
	IsBasic         bool                        `json:"is_basic" tf:"is_basic"`
	Description     string                      `json:"description" tf:"description"`
	Expression      string                      `json:"expression" tf:"expression"`
	BasicExpression []*SuppressionRuleCondition `json:"basic_expression" tf:"basic_expressions"`
	IsTimeBased     bool                        `json:"is_timebased" tf:"is_timebased"`
	TimeSlots       []*TimeSlot                 `json:"timeslots" tf:"timeslots"`
}

type TimeSlot struct {
	TimeZone   string      `json:"time_zone" tf:"time_zone"`
	StartTime  string      `json:"start_time" tf:"start_time"`
	EndTime    string      `json:"end_time" tf:"end_time"`
	IsAllDay   bool        `json:"is_allday" tf:"is_allday"`
	Repetition string      `json:"repetition" tf:"repetition"`
	IsCustom   bool        `json:"is_custom" tf:"is_custom"`
	Custom     *CustomTime `json:"custom" tf:"custom"`
	EndsNever  bool        `json:"ends_never" tf:"ends_never"`
	EndsOn     string      `json:"ends_on" tf:"ends_on"`
}

func (c *TimeSlot) Encode() (tf.M, error) {
	return tf.Encode(c)
}

type CustomTime struct {
	RepeatsCount      int    `json:"repeats_count" tf:"repeats_count"`
	Repeats           string `json:"repeats" tf:"repeats"`
	RepeatsOnWeekdays []int  `json:"repeats_on_weekdays" tf:"repeats_on_weekdays"`
	RepeatsOnMonth    string `json:"repeats_on_month" tf:"repeats_on_month"`
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

	rtimeSlots := r.TimeSlots

	timeSlots, err := tf.EncodeSlice(r.TimeSlots)
	if err != nil {
		return nil, err
	}
	m["timeslots"] = timeSlots

	if rtimeSlots == nil {
		rtimeSlots = []*TimeSlot{}
		m["is_timebased"] = false
	} else {
		for idx, t := range rtimeSlots {
			mNewCustomField := tf.List(tf.M{})
			if t.Repetition == "custom" {
				mNewCustomField = tf.List(tf.M{
					"repeats_count":       t.Custom.RepeatsCount,
					"repeats":             t.Custom.Repeats,
					"repeats_on_weekdays": t.Custom.RepeatsOnWeekdays,
					"repeats_on_month":    t.Custom.RepeatsOnMonth,
				})
			}
			m["timeslots"].([]interface{})[idx].(map[string]interface{})["custom"] = mNewCustomField
		}
	}

	return m, nil
}

type SuppressionRules struct {
	ID        string             `json:"id" tf:"id"`
	ServiceID string             `json:"service_id" tf:"service_id"`
	Rules     []*SuppressionRule `json:"rules" tf:"-"`
	Rule      *SuppressionRule   `json:"rule" tf:"-"`
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
	url := fmt.Sprintf("%s/services/%s/suppression-rules?owner_id=%s", client.BaseURLV3, serviceID, teamID)

	return Request[any, SuppressionRules](http.MethodGet, url, client, ctx, nil)
}

type UpdateSuppressionRulesReq struct {
	Rules []SuppressionRule `json:"rules"`
}

func (client *Client) UpdateSuppressionRules(ctx context.Context, serviceID, teamID string, req *UpdateSuppressionRulesReq) (*SuppressionRules, error) {
	url := fmt.Sprintf("%s/services/%s/suppression-rules", client.BaseURLV3, serviceID)
	return Request[UpdateSuppressionRulesReq, SuppressionRules](http.MethodPost, url, client, ctx, req)
}

// suppression rules v2

type CreateSuppressionRule struct {
	Rule SuppressionRule `json:"rule"`
}

func (client *Client) CreateSuppressionRulesV2(ctx context.Context, serviceID string, req *CreateSuppressionRule) (*CreateSuppressionRule, error) {
	url := fmt.Sprintf("%s/services/%s/suppression-rules/new", client.BaseURLV3, serviceID)
	return Request[CreateSuppressionRule, CreateSuppressionRule](http.MethodPost, url, client, ctx, req)
}

func (client *Client) UpdateSuppressionRuleByID(ctx context.Context, serviceID, ruleID string, req *CreateSuppressionRule) (*CreateSuppressionRule, error) {
	url := fmt.Sprintf("%s/services/%s/suppression-rules/%s", client.BaseURLV3, serviceID, ruleID)
	return Request[CreateSuppressionRule, CreateSuppressionRule](http.MethodPut, url, client, ctx, req)
}

type SuppressionRuleV2 struct {
	ID        string           `json:"rule_id" tf:"id"`
	ServiceID string           `json:"service_id" tf:"service_id"`
	Rule      *SuppressionRule `json:"rule" tf:"-"`
}

func (s *SuppressionRuleV2) Encode() (tf.M, error) {
	m, err := tf.Encode(s.Rule)
	if err != nil {
		return nil, err
	}

	basicExpressions, err := tf.EncodeSlice(s.Rule.BasicExpression)
	if err != nil {
		return nil, err
	}
	m["basic_expressions"] = basicExpressions

	rtimeSlots := s.Rule.TimeSlots

	timeSlots, err := tf.EncodeSlice(rtimeSlots)
	if err != nil {
		return nil, err
	}
	m["timeslots"] = timeSlots

	if len(rtimeSlots) > 0 {
		for idx, t := range rtimeSlots {
			m["timeslots"].([]interface{})[idx].(map[string]interface{})["custom"] = nil
			if t.Repetition == "custom" {
				mNewCustomField := tf.List(tf.M{
					"repeats_count":       t.Custom.RepeatsCount,
					"repeats":             t.Custom.Repeats,
					"repeats_on_weekdays": t.Custom.RepeatsOnWeekdays,
					"repeats_on_month":    t.Custom.RepeatsOnMonth,
				})
				m["timeslots"].([]interface{})[idx].(map[string]interface{})["custom"] = mNewCustomField
			}
		}
	}

	return m, nil
}

func (client *Client) GetSuppressionRuleByID(ctx context.Context, serviceID, ruleID string) (*SuppressionRuleV2, error) {
	url := fmt.Sprintf("%s/services/%s/suppression-rules/%s", client.BaseURLV3, serviceID, ruleID)

	return Request[any, SuppressionRuleV2](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) DeleteSuppressionRuleByID(ctx context.Context, serviceID, ruleID string) (any, error) {
	url := fmt.Sprintf("%s/services/%s/suppression-rules/%s", client.BaseURLV3, serviceID, ruleID)

	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
