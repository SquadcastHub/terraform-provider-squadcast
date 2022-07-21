package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type Data struct {
	Slo *Slo `json:"slo,omitempty"`
}

type Slo struct {
	ID                  uint                  `json:"id,omitempty" tf:"id"`
	Name                string                `json:"name" tf:"name"`
	Description         string                `json:"description,omitempty" tf:"description"`
	TimeIntervalType    string                `json:"time_interval_type" tf:"time_interval_type"`
	ServiceIDs          []string              `json:"service_ids" tf:"service_ids"`
	Slis                []string              `json:"slis" tf:"slis"`
	TargetSlo           float64               `json:"target_slo" tf:"target_slo"`
	StartTime           string                `json:"start_time,omitempty" tf:"start_time"`
	EndTime             string                `json:"end_time,omitempty" tf:"end_time"`
	DurationInDays      int                   `json:"duration_in_days,omitempty" tf:"duration_in_days"`
	SloMonitoringChecks []*SloMonitoringCheck `json:"slo_monitoring_checks" tf:"rules"`
	SloActions          []*SloAction          `json:"slo_actions" tf:"notify"`
	OwnerID             string                `json:"owner_id" tf:"team_id"`
}

type SloMonitoringCheck struct {
	ID        uint   `json:"id,omitempty" tf:"id"`
	SloID     int64  `json:"slo_id,omitempty" tf:"slo_id"`
	Name      string `json:"name" tf:"name"`
	Threshold int    `json:"threshold" tf:"threshold"`
	IsChecked bool   `json:"is_checked" tf:"is_checked"`
}

type SloAction struct {
	ID        uint   `json:"id,omitempty" tf:"id"`
	SloID     int64  `json:"slo_id,omitempty" tf:"slo_id"`
	Type      string `json:"type" tf:"type"`
	UserID    string `json:"user_id" tf:"user_id"`
	SquadID   string `json:"squad_id" tf:"squad_id"`
	ServiceID string `json:"service_id" tf:"service_id"`
}

type SloNotify struct {
	ID        uint     `json:"id,omitempty" tf:"id"`
	SloID     int64    `json:"slo_id,omitempty" tf:"slo_id"`
	UserIDs   []string `json:"user_ids" tf:"user_ids"`
	SquadIDs  []string `json:"squad_ids" tf:"squad_ids"`
	ServiceID string   `json:"service_id" tf:"service_id"`
}

func (c *SloMonitoringCheck) Encode() (map[string]interface{}, error) {
	return tf.Encode(c)
}

func (c *SloNotify) Encode() (map[string]interface{}, error) {
	return tf.Encode(c)
}

func (r *Slo) Encode() (map[string]interface{}, error) {
	notify := make([]*SloNotify, 0)
	// Max item limit set to 1 for `notify` rule,
	// However, notify type is still a list. so we are taking slice
	notify = append(notify, &SloNotify{})

	m, err := tf.Encode(r)
	if err != nil {
		return nil, err
	}

	sloMonitoringChecks, err := tf.EncodeSlice(r.SloMonitoringChecks)
	if err != nil {
		return nil, err
	}
	m["rules"] = sloMonitoringChecks

	for _, n := range r.SloActions {
		if n.UserID != "" {
			notify[0].UserIDs = append(notify[0].UserIDs, n.UserID)
		}
		if n.SquadID != "" {
			notify[0].SquadIDs = append(notify[0].SquadIDs, n.SquadID)
		}
		if n.ServiceID != "" {
			notify[0].ServiceID = n.ServiceID
		}
	}

	if len(r.SloActions) > 0 {
		notify[0].SloID = int64(r.ID)
	}

	notifyObj, err := tf.EncodeSlice(notify)
	if err != nil {
		fmt.Println(err)

	}
	m["notify"] = notifyObj

	return m, nil
}

func (r *Data) Encode() (map[string]interface{}, error) {
	m, err := tf.Encode(r)
	if err != nil {
		return nil, err
	}

	slo, err := tf.Encode(r.Slo)
	if err != nil {
		return nil, err
	}
	m["slo"] = slo

	return m, nil
}

func (client *Client) CreateSlo(ctx context.Context, orgID, ownerID string, req *Slo) (*Slo, error) {
	url := fmt.Sprintf("%s/slo?owner_type=team&owner_id=%s", client.BaseURLV3, ownerID)
	data, err := Request[Slo, Data](http.MethodPost, url, client, ctx, req)
	if err != nil {
		return nil, err
	}
	return data.Slo, err
}

func (client *Client) GetSlo(ctx context.Context, orgID, ownerID, sloID string) (*Slo, error) {
	url := fmt.Sprintf("%s/slo/%s?owner_type=team&owner_id=%s", client.BaseURLV3, sloID, ownerID)
	data, err := Request[any, Data](http.MethodGet, url, client, ctx, nil)
	if err != nil {
		return nil, errors.New("Slo not found")
	}
	return data.Slo, nil
}

func (client *Client) UpdateSlo(ctx context.Context, orgID, ownerID, sloID string, req *Slo) (*Slo, error) {
	url := fmt.Sprintf("%s/slo/%s?owner_type=team&owner_id=%s", client.BaseURLV3, sloID, ownerID)
	return Request[Slo, Slo](http.MethodPut, url, client, ctx, req)
}

func (client *Client) DeleteSlo(ctx context.Context, orgID, ownerID, sloID string) (*any, error) {
	url := fmt.Sprintf("%s/slo/%s?owner_type=team&owner_id=%s", client.BaseURLV3, sloID, ownerID)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
