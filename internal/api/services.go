package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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
	ActiveAlertSources []string           `json:"-" tf:"alert_sources"`
	AlertSources       map[string]string  `json:"-" tf:"alert_source_endpoints"`
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
