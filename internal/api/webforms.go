package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type WebformReq struct {
	TeamID        string            `json:"owner_id"`
	Name          string            `json:"name"`
	IsCname       bool              `json:"is_cname"`
	PublicUrl     string            `json:"public_url"`
	HostName      string            `json:"host_name"`
	Tags          map[string]string `json:"tags"`
	FormOwnerType string            `json:"form_owner_type"`
	FormOwnerID   string            `json:"form_owner_id"`
	Services      []WFService       `json:"services"`
	InputField    []WFInputField    `json:"input_field"`
	Header        string            `json:"header"`
	Title         string            `json:"title"`
	FooterText    string            `json:"footer_text"`
	FooterLink    string            `json:"footer_link"`
	EmailOn       []string          `json:"email_on"`
	Description   string            `json:"description"`
}

type Webform struct {
	ID            uint              `json:"id" tf:"id"`
	TeamID        string            `json:"owner_id" tf:"team_id"`
	Name          string            `json:"name" tf:"name"`
	PublicUrl     string            `json:"public_url" tf:"public_url"`
	HostName      string            `json:"host_name" tf:"custom_domain_name"`
	Tags          map[string]string `json:"tags" tf:"tags"`
	FormOwnerType string            `json:"form_owner_type"`
	FormOwnerID   string            `json:"form_owner_id"`
	FormOwnerName string            `json:"form_owner_name"`
	WebformOwner  *WebformOwner     `tf:"owner"`
	Services      []WFService       `json:"services" tf:"services"`
	InputField    []WFInputField    `json:"input_field" tf:"input_field"`
	Header        string            `json:"header" tf:"header"`
	Title         string            `json:"title" tf:"title"`
	FooterText    string            `json:"footer_text" tf:"footer_text"`
	FooterLink    string            `json:"footer_link" tf:"footer_link"`
	EmailOn       []string          `json:"email_on" tf:"email_on"`
	Description   string            `json:"description" tf:"description"`
}

type CreateWebformRes struct {
	WebFormRes *Webform `json:"webform"`
}

type WFService struct {
	ServiceId string `json:"service_id" tf:"service_id"`
	Name      string `json:"name" tf:"name"`
	Alias     string `json:"alias" tf:"alias"`
}

type WFTag struct {
	Key   string `json:"key" tf:"key"`
	Value string `json:"value" tf:"value"`
}

type WFInputField struct {
	Label   string   `json:"label" tf:"label"`
	Options []string `json:"options" tf:"options"`
}

type WebformOwner struct {
	ID   string `tf:"id"`
	Name string `tf:"name"`
	Type string `tf:"type"`
}

func (webformTag WFTag) Encode() (tf.M, error) {
	return tf.Encode(webformTag)
}

func (webformService WFService) Encode() (tf.M, error) {
	return tf.Encode(webformService)
}

func (webformInputField WFInputField) Encode() (tf.M, error) {
	return tf.Encode(webformInputField)
}

func (t *Webform) Encode() (tf.M, error) {
	m, err := tf.Encode(t)
	if err != nil {
		return nil, err
	}
	m["team_id"] = t.TeamID

	m["owner"] = tf.List(tf.M{
		"id":   t.FormOwnerID,
		"name": t.FormOwnerName,
		"type": t.FormOwnerType,
	})

	m["custom_domain_name"] = t.HostName

	tags, err := tf.Encode(t.Tags)
	if err != nil {
		return nil, err
	}
	m["tags"] = tags

	services, err := tf.EncodeSlice(t.Services)
	if err != nil {
		return nil, err
	}
	m["services"] = services

	inputFieldEncoded, err := tf.EncodeSlice(t.InputField)
	if err != nil {
		return nil, err
	}
	m["input_field"] = inputFieldEncoded

	return m, nil
}

func (client *Client) GetWebformById(ctx context.Context, teamID string, id string) (*Webform, error) {
	url := fmt.Sprintf("%s/webform/%s?owner_id=%s", client.BaseURLV3, id, teamID)

	return Request[any, Webform](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) GetWebformByName(ctx context.Context, teamID string, name string) (*Webform, error) {
	url := fmt.Sprintf("%s/webform/by-name?name=%s&owner_id=%s", client.BaseURLV3, url.QueryEscape(name), teamID)

	return Request[any, Webform](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) CreateWebform(ctx context.Context, teamID string, req *WebformReq) (*CreateWebformRes, error) {
	url := fmt.Sprintf("%s/webform?owner_id=%s", client.BaseURLV3, teamID)

	return Request[WebformReq, CreateWebformRes](http.MethodPost, url, client, ctx, req)
}

func (client *Client) UpdateWebform(ctx context.Context, teamID string, id string, req *WebformReq) (*Webform, error) {
	url := fmt.Sprintf("%s/webform/%s?owner_id=%s", client.BaseURLV3, id, teamID)

	return Request[WebformReq, Webform](http.MethodPut, url, client, ctx, req)
}

func (client *Client) DeleteWebform(ctx context.Context, teamID string, id string) (*any, error) {
	url := fmt.Sprintf("%s/webform/%s?owner_id=%s", client.BaseURLV3, id, teamID)

	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
