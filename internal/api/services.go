package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type Service struct {
	ID                 string             `json:"id" tf:"id"`
	Name               string             `json:"name" tf:"name"`
	APIKey             string             `json:"api_key" tf:"api_key"`
	Email              string             `json:"email" tf:"email"`
	EmailPrefix        string             `json:"-" tf:"email_prefix"`
	Description        string             `json:"description" tf:"description"`
	EscalationPolicyID string             `json:"escalation_policy_id" tf:"escalation_policy_id"`
	OnMaintenance      bool               `json:"on_maintenance" tf:"-"`
	Owner              OwnerRef           `json:"owner" tf:"-"`
	Maintainer         *ServiceMaintainer `json:"maintainer" tf:"maintainer"`
	Tags               []ServiceTag       `json:"tags" tf:"tags"`
	Dependencies       []string           `json:"depends" tf:"dependencies"`
	ActiveAlertSources map[string]string  `json:"-" tf:"active_alert_source_endpoints"`
	AlertSources       map[string]string  `json:"-" tf:"alert_source_endpoints"`
	Slack              *SlackChannel      `json:"slack" tf:"-"`

	APTAConfig              APTAConfig               `json:"auto_pause_transient_alerts_config" tf:"-"`
	IAGConfig               IAGConfig                `json:"intelligent_alerts_grouping_config" tf:"-"`
	DelayNotificationConfig NotificationsDelayConfig `json:"delay_notification_config" tf:"-"`
}

type SlackChannel struct {
	ChannelID string `json:"channel_id" tf:"-"`
}

type APTAConfig struct {
	IsEnabled     bool `json:"is_enabled" tf:"is_enabled"`
	TimeoutInMins int  `json:"timeout_in_mins" tf:"timeout"`
}

type IAGConfig struct {
	IsEnabled           bool `json:"is_enabled" tf:"is_enabled"`
	RollingWindowInMins int  `json:"rolling_window_in_mins" tf:"grouping_window"`
}

type NotificationsDelayConfig struct {
	IsEnabled              bool                       `json:"is_enabled" tf:"is_enabled"`
	Timezone               string                     `json:"timezone" tf:"timezone"`
	FixedTimeslotConfig    *FixedTimeslotConfig       `json:"fixed_timeslot_config,omitempty" tf:"fixed_timeslot_config"`
	CustomTimeslotsEnabled bool                       `json:"custom_timeslots_enabled" tf:"custom_timeslots_enabled"`
	CustomTimeslots        map[string][]DelayTimeSlot `json:"custom_timeslots,omitempty" tf:"custom_timeslots"`
	AssignedTo             *AssignTo                  `json:"assigned_to" tf:"assigned_to"`
}
type AssignTo struct {
	ID   string `json:"id,omitempty" tf:"id"`
	Type string `json:"type,omitempty" tf:"type"`
}
type FixedTimeslotConfig struct {
	DelayTimeSlot
	RepeatOnDays []int `json:"repeat_days,omitempty" tf:"repeat_days"`
}
type DelayTimeSlot struct {
	StartTime string `json:"start_time,omitempty" tf:"start_time"`
	EndTime   string `json:"end_time,omitempty" tf:"end_time"`
}

var days = []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}
var DayOfWeekMap = map[string]time.Weekday{
	"sunday":    time.Sunday,
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
}

func (apta APTAConfig) Encode() (tf.M, error) {
	m, err := tf.Encode(apta)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (iag IAGConfig) Encode() (tf.M, error) {
	m, err := tf.Encode(iag)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (cfg NotificationsDelayConfig) Encode() (tf.M, error) {
	delayConfig := tf.M{
		"is_enabled":               cfg.IsEnabled,
		"timezone":                 cfg.Timezone,
		"custom_timeslots_enabled": cfg.CustomTimeslotsEnabled,
	}
	if cfg.FixedTimeslotConfig != nil {
		if cfg.FixedTimeslotConfig.StartTime != "" && cfg.FixedTimeslotConfig.EndTime != "" && len(cfg.FixedTimeslotConfig.RepeatOnDays) > 0 {
			repeatDays := make([]string, 0, len(cfg.FixedTimeslotConfig.RepeatOnDays))
			for _, repeatDay := range cfg.FixedTimeslotConfig.RepeatOnDays {
				repeatDays = append(repeatDays, strings.ToLower(time.Weekday(repeatDay).String()))
			}
			delayConfig["fixed_timeslot_config"] = tf.List(tf.M{
				"start_time":  cfg.FixedTimeslotConfig.StartTime,
				"end_time":    cfg.FixedTimeslotConfig.EndTime,
				"repeat_days": repeatDays,
			})
		}
	}
	if cfg.AssignedTo != nil {
		delayConfig["assigned_to"] = tf.List(tf.M{
			"id":   cfg.AssignedTo.ID,
			"type": cfg.AssignedTo.Type,
		})
	}
	if cfg.CustomTimeslots != nil {
		customTimeSlots := make([]tf.M, 0, len(cfg.CustomTimeslots))
		for k, v := range cfg.CustomTimeslots {
			for _, slot := range v {
				day, e := strconv.Atoi(k)
				if e != nil {
					return nil, e
				}
				customTimeSlots = append(customTimeSlots, tf.M{
					"day_of_week": days[day],
					"start_time":  slot.StartTime,
					"end_time":    slot.EndTime,
				})
			}
		}
		delayConfig["custom_timeslots"] = customTimeSlots
	}

	return delayConfig, nil
}

func (serviceTag ServiceTag) Encode() (tf.M, error) {
	return tf.Encode(serviceTag)
}

func (s *Service) Encode() (tf.M, error) {
	s.EmailPrefix = strings.Split(s.Email, "@")[0]

	m, err := tf.Encode(s)
	if err != nil {
		return nil, err
	}

	m["team_id"] = s.Owner.ID

	if s.Maintainer != nil {
		m["maintainer"] = tf.List(tf.M{
			"type": s.Maintainer.Type,
			"id":   s.Maintainer.ID,
		})
	}

	tagsEncoded, terr := tf.EncodeSlice(s.Tags)
	if terr != nil {
		return nil, terr
	}
	m["tags"] = tagsEncoded

	if s.Slack != nil {
		m["slack_channel_id"] = s.Slack.ChannelID
	}

	return m, nil
}

func (client *Client) GetServiceById(ctx context.Context, teamID string, id string) (*Service, error) {
	url := fmt.Sprintf("%s/services/%s?owner_id=%s", client.BaseURLV3, id, teamID)

	return Request[any, Service](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) GetServiceByName(ctx context.Context, teamID string, name string) (*Service, error) {
	url := fmt.Sprintf("%s/services/by-name?name=%s&owner_id=%s", client.BaseURLV3, url.QueryEscape(name), teamID)

	return Request[any, Service](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) ListServices(ctx context.Context, teamID string) ([]*Service, error) {
	url := fmt.Sprintf("%s/services?owner_id=%s", client.BaseURLV3, teamID)

	return RequestSlice[any, Service](http.MethodGet, url, client, ctx, nil)
}

type CreateServiceReq struct {
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	TeamID             string             `json:"owner_id"`
	EscalationPolicyID string             `json:"escalation_policy_id"`
	EmailPrefix        string             `json:"email_prefix"`
	Maintainer         *ServiceMaintainer `json:"maintainer"`
	Tags               []ServiceTag       `json:"tags"`
}

type UpdateServiceReq struct {
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	EscalationPolicyID string             `json:"escalation_policy_id"`
	EmailPrefix        string             `json:"email_prefix"`
	Maintainer         *ServiceMaintainer `json:"maintainer"`
	Tags               []ServiceTag       `json:"tags"`
}

type ServiceMaintainer struct {
	ID   string `json:"id" tf:"id"`
	Type string `json:"type" tf:"type"`
}

type ServiceTag struct {
	Key   string `json:"key" tf:"key"`
	Value string `json:"value" tf:"value"`
}

type AddSlackChannelReq struct {
	ChannelID string `json:"channel_id"`
}

type UpdateServiceDependenciesReq struct {
	Data []string `json:"data"`
}

func (client *Client) CreateService(ctx context.Context, req *CreateServiceReq) (*Service, error) {
	url := fmt.Sprintf("%s/services", client.BaseURLV3)
	return Request[CreateServiceReq, Service](http.MethodPost, url, client, ctx, req)
}

func (client *Client) UpdateService(ctx context.Context, id string, req *UpdateServiceReq) (*Service, error) {
	url := fmt.Sprintf("%s/services/%s", client.BaseURLV3, id)
	return Request[UpdateServiceReq, Service](http.MethodPut, url, client, ctx, req)
}

func (client *Client) UpdateServiceDependencies(ctx context.Context, id string, req *UpdateServiceDependenciesReq) (*any, error) {
	url := fmt.Sprintf("%s/organizations/%s/services/%s/dependencies", client.BaseURLV2, client.OrganizationID, id)
	return Request[UpdateServiceDependenciesReq, any](http.MethodPost, url, client, ctx, req)
}

func (client *Client) DeleteService(ctx context.Context, id string) (*any, error) {
	url := fmt.Sprintf("%s/services/%s", client.BaseURLV3, id)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}

func (client *Client) UpdateSlackChannel(ctx context.Context, serviceID string, req *AddSlackChannelReq) (*any, error) {
	url := fmt.Sprintf("%s/services/%s/extensions", client.BaseURLV3, serviceID)
	return Request[AddSlackChannelReq, any](http.MethodPut, url, client, ctx, req)
}

// IAG, APTA, Delayed Notification
func (client *Client) UpdateAPTAConfig(ctx context.Context, serviceID string, req *APTAConfig) (*any, error) {
	url := fmt.Sprintf("%s/services/%s/apta-config", client.BaseURLV3, serviceID)
	return Request[APTAConfig, any](http.MethodPut, url, client, ctx, req)
}

func (client *Client) UpdateIAGConfig(ctx context.Context, serviceID string, req *IAGConfig) (*any, error) {
	url := fmt.Sprintf("%s/services/%s/iag-config", client.BaseURLV3, serviceID)
	return Request[IAGConfig, any](http.MethodPut, url, client, ctx, req)
}

func (client *Client) UpdateDelayedNotificationConfig(ctx context.Context, serviceID string, req *NotificationsDelayConfig) (*any, error) {
	url := fmt.Sprintf("%s/services/%s/notification-delay-config", client.BaseURLV3, serviceID)
	return Request[NotificationsDelayConfig, any](http.MethodPut, url, client, ctx, req)
}
