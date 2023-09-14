package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type StatusPage struct {
	ID                           uint   `json:"id,omitempty" tf:"id"`
	TeamID                       string `json:"teamID" tf:"team_id"`
	Name                         string `json:"name" tf:"name"`
	Description                  string `json:"description,omitempty" tf:"description"`
	IsPublic                     bool   `json:"isPublic" tf:"is_public"`
	DomainName                   string `json:"domainName" tf:"domain_name"`
	CustomDomainName             string `json:"customDomainName,omitempty" tf:"custom_domain_name"`
	Timezone                     string `json:"timezone" tf:"timezone"`
	ContactEmail                 string `json:"contactEmail" tf:"contact_email"`
	ThemeColor                   `json:"themeColor" tf:"theme_color"`
	OwnerType                    string           `json:"ownerType"`
	OwnerID                      string           `json:"ownerID"`
	StatusPageOwner              *StatusPageOwner `tf:"owner"`
	AllowWebhookSubscription     bool             `json:"allowWebhookSubscription" tf:"allow_webhook_subscription"`
	AllowMaintenanceSubscription bool             `json:"allowMaintenanceSubscription" tf:"allow_maintenance_subscription"`
	AllowComponentsSubscription  bool             `json:"allowComponentsSubscription" tf:"allow_components_subscription"`
}

type StatusPageOwner struct {
	ID   string `tf:"id"`
	Type string `tf:"type"`
}

type ThemeColor struct {
	Primary   string `json:"primary" tf:"primary"`
	Secondary string `json:"secondary" tf:"secondary"`
}

type StatusPageComponent struct {
	ID             uint   `json:"id,omitempty" tf:"id"`
	PageID         uint   `json:"pageID" tf:"status_page_id"`
	Name           string `json:"name" tf:"name"`
	Description    string `json:"description,omitempty" tf:"description"`
	GroupID        *uint  `json:"groupID,omitempty" tf:"group_id"`
	BelongsToGroup *bool  `json:"belongsToGroup" tf:"-"`
}

type StatusPageGroup struct {
	ID     uint   `json:"id,omitempty" tf:"id"`
	PageID uint   `json:"pageID" tf:"status_page_id"`
	Name   string `json:"name" tf:"name"`
}

func (sp *StatusPage) Encode() (map[string]interface{}, error) {
	m, err := tf.Encode(sp)
	if err != nil {
		return nil, err
	}

	m["theme_color"] = tf.List(tf.M{
		"primary":   sp.ThemeColor.Primary,
		"secondary": sp.ThemeColor.Secondary,
	})

	m["owner"] = tf.List(tf.M{
		"id":   sp.OwnerID,
		"type": sp.OwnerType,
	})

	return m, nil
}

func (spc *StatusPageComponent) Encode() (map[string]interface{}, error) {
	m, err := tf.Encode(spc)
	if err != nil {
		return nil, err
	}

	statusPageID := strconv.FormatUint(uint64(spc.PageID), 10)
	m["status_page_id"] = statusPageID

	if spc.GroupID != nil {
		groupID := strconv.FormatUint(uint64(*spc.GroupID), 10)
		m["group_id"] = groupID
	}

	return m, nil
}

func (spg *StatusPageGroup) Encode() (map[string]interface{}, error) {
	m, err := tf.Encode(spg)
	if err != nil {
		return nil, err
	}

	statusPageID := strconv.FormatUint(uint64(spg.PageID), 10)
	m["status_page_id"] = statusPageID

	return m, nil
}

func (client *Client) CreateStatusPage(ctx context.Context, req *StatusPage) (*StatusPage, error) {
	url := fmt.Sprintf("%s/statuspages", client.BaseURLV4)
	data, err := Request[StatusPage, StatusPage](http.MethodPost, url, client, ctx, req)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (client *Client) GetStatusPageById(ctx context.Context, ID string) (*StatusPage, error) {
	url := fmt.Sprintf("%s/statuspages/%s", client.BaseURLV4, ID)
	data, err := Request[any, StatusPage](http.MethodGet, url, client, ctx, nil)
	if err != nil {
		return nil, errors.New("Status page with given ID not found")
	}
	return data, nil
}

func (client *Client) UpdateStatusPage(ctx context.Context, pageID string, req *StatusPage) (*StatusPage, error) {
	url := fmt.Sprintf("%s/statuspages/%s", client.BaseURLV4, pageID)
	return Request[StatusPage, StatusPage](http.MethodPut, url, client, ctx, req)
}

func (client *Client) DeleteStatusPage(ctx context.Context, ID string) (*any, error) {
	url := fmt.Sprintf("%s/statuspages/%s", client.BaseURLV4, ID)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}

// Status Page Component APIs
func (client *Client) CreateStatusPageComponent(ctx context.Context, pageID string, req *StatusPageComponent) (*StatusPageComponent, error) {
	url := fmt.Sprintf("%s/statuspages/%s/components", client.BaseURLV4, pageID)
	data, err := Request[StatusPageComponent, StatusPageComponent](http.MethodPost, url, client, ctx, req)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (client *Client) GetStatusPageComponentById(ctx context.Context, pageID, componentID string) (*StatusPageComponent, error) {
	url := fmt.Sprintf("%s/statuspages/%s/components/%s", client.BaseURLV4, pageID, componentID)
	data, err := Request[any, StatusPageComponent](http.MethodGet, url, client, ctx, nil)
	if err != nil {
		return nil, errors.New("Status page with given ID not found")
	}
	return data, nil
}

func (client *Client) UpdateStatusPageComponent(ctx context.Context, pageID, componentID string, req *StatusPageComponent) (*StatusPageComponent, error) {
	url := fmt.Sprintf("%s/statuspages/%s/components/%s", client.BaseURLV4, pageID, componentID)
	return Request[StatusPageComponent, StatusPageComponent](http.MethodPut, url, client, ctx, req)
}

func (client *Client) DeleteStatusPageComponent(ctx context.Context, pageID, componentID string) (*any, error) {
	url := fmt.Sprintf("%s/statuspages/%s/components/%s", client.BaseURLV4, pageID, componentID)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}

// Status Page Group APIs
func (client *Client) CreateStatusPageGroup(ctx context.Context, pageID string, req *StatusPageGroup) (*StatusPageGroup, error) {
	url := fmt.Sprintf("%s/statuspages/%s/groups", client.BaseURLV4, pageID)
	data, err := Request[StatusPageGroup, StatusPageGroup](http.MethodPost, url, client, ctx, req)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (client *Client) GetStatusPageGroupById(ctx context.Context, pageID, groupID string) (*StatusPageGroup, error) {
	url := fmt.Sprintf("%s/statuspages/%s/groups/%s", client.BaseURLV4, pageID, groupID)
	data, err := Request[any, StatusPageGroup](http.MethodGet, url, client, ctx, nil)
	if err != nil {
		return nil, errors.New("Status page with given ID not found")
	}
	return data, nil
}

func (client *Client) UpdateStatusPageGroup(ctx context.Context, pageID, groupID string, req *StatusPageGroup) (*StatusPageGroup, error) {
	url := fmt.Sprintf("%s/statuspages/%s/groups/%s", client.BaseURLV4, pageID, groupID)
	return Request[StatusPageGroup, StatusPageGroup](http.MethodPut, url, client, ctx, req)
}

func (client *Client) DeleteStatusPageGroup(ctx context.Context, pageID, groupID string) (*any, error) {
	url := fmt.Sprintf("%s/statuspages/%s/groups/%s", client.BaseURLV4, pageID, groupID)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
