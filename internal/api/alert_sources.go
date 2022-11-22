package api

import (
	"context"
	"fmt"
	"net/http"
)

type AlertSource struct {
	ID             string `json:"_id"`
	Type           string `json:"type"`
	Heading        string `json:"heading"`
	SupportDocURL  string `json:"supportDoc"`
	DisplayKeyOnly bool   `json:"displayKeyOnly"`
	ShortName      string `json:"shortName"`
	Version        string `json:"version"`

	IsValid      bool `json:"isValid"`
	IsPrivate    bool `json:"isPrivate"`
	IsDeprecated bool `json:"deprecated"`
}

type AlertSourcesList []*AlertSource

type ActiveAlertSource struct {
	ID string `json:"alert_source_id"`
}

type ActiveAlertSources struct {
	AlertSources []ActiveAlertSource `json:"alert_sources"`
}

type AddAlertSourcesReq struct {
	AlertSources []string `json:"alert_sources"`
}

func (asl *AlertSourcesList) Available() *AlertSourcesList {
	var list AlertSourcesList

	for _, v := range *asl {
		if v.IsValid && !v.IsPrivate && !v.IsDeprecated {
			list = append(list, v)
		}
	}

	return &list
}

func (asl *AlertSourcesList) EndpointMap(ingestionBaseURL string, service *Service) map[string]string {
	m := make(map[string]string, len(*asl))

	for _, as := range *asl {
		m[as.ShortName] = as.Endpoint(ingestionBaseURL, service)
	}

	return m
}

func (alertSource *AlertSource) Endpoint(ingestionBaseURL string, service *Service) string {
	if alertSource.ShortName == "email" {
		return service.Email
	}

	if alertSource.DisplayKeyOnly {
		return service.APIKey
	}

	return fmt.Sprintf("%s/%s/incidents/%s/%s", ingestionBaseURL, alertSource.Version, alertSource.ShortName, service.APIKey)
}

func (client *Client) ListAlertSources(ctx context.Context) (AlertSourcesList, error) {
	url := fmt.Sprintf("%s/public/integrations", client.BaseURLV2)

	return RequestSlice[any, AlertSource](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) AddAlertSources(ctx context.Context, serviceID string, alertSources *AddAlertSourcesReq) (*any, error) {
	url := fmt.Sprintf("%s/catalog-services/%s/alert-sources", client.BaseURLV3, serviceID)
	return Request[AddAlertSourcesReq, any](http.MethodPut, url, client, ctx, alertSources)
}

func (client *Client) ListActiveAlertSources(ctx context.Context, serviceID string) (*ActiveAlertSources, error) {
	url := fmt.Sprintf("%s/catalog-services/%s/alert-sources", client.BaseURLV3, serviceID)
	return Request[any, ActiveAlertSources](http.MethodGet, url, client, ctx, nil)
}
