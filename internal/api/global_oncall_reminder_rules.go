package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type GlobalOncallReminderRules struct {
	ID int `json:"id" tf:"id"`
	GlobalOncallReminderRulesReq
}

type GlobalOncallReminderRulesReq struct {
	IsEnabled bool                `json:"is_enabled" tf:"is_enabled"`
	TeamID    string              `json:"owner_id,omitempty" tf:"team_id"`
	Rules     []*NotificationRule `json:"rules" tf:"rules"`
}

type NotificationRule struct {
	TypeOfNotification string `json:"type" tf:"type"`
	Time               int    `json:"time" tf:"time"`
}

func (c NotificationRule) Encode() (tf.M, error) {
	return tf.Encode(c)
}

func (g *GlobalOncallReminderRules) Encode() (tf.M, error) {
	m, err := tf.Encode(g)
	if err != nil {
		return nil, err
	}
	m["id"] = strconv.Itoa(g.ID)

	rules, err := tf.EncodeSlice(g.Rules)
	if err != nil {
		return nil, err
	}
	m["rules"] = rules

	return m, nil
}

func (client *Client) CreateGlobalOncallReminderRules(ctx context.Context, req *GlobalOncallReminderRulesReq) (*GlobalOncallReminderRules, error) {
	url := fmt.Sprintf("%s/global-oncall-reminder-rules", client.BaseURLV3)

	return Request[GlobalOncallReminderRulesReq, GlobalOncallReminderRules](http.MethodPost, url, client, ctx, req)
}

func (client *Client) UpdateGlobalOncallReminderRules(ctx context.Context, teamID string, req *GlobalOncallReminderRulesReq) (*GlobalOncallReminderRules, error) {
	url := fmt.Sprintf("%s/global-oncall-reminder-rules?owner_id=%s", client.BaseURLV3, teamID)

	return Request[GlobalOncallReminderRulesReq, GlobalOncallReminderRules](http.MethodPut, url, client, ctx, req)
}

func (client *Client) DeleteGlobalOncallReminderRules(ctx context.Context, teamID string) (*any, error) {
	url := fmt.Sprintf("%s/global-oncall-reminder-rules?owner_id=%s", client.BaseURLV3, teamID)

	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}

func (client *Client) GetGlobalOncallReminderRules(ctx context.Context, teamID string) (*GlobalOncallReminderRules, error) {
	url := fmt.Sprintf("%s/global-oncall-reminder-rules?owner_id=%s", client.BaseURLV3, teamID)

	return Request[any, GlobalOncallReminderRules](http.MethodGet, url, client, ctx, nil)
}
