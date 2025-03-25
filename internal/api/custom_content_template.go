package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type CustomContentTemplateOverlayReq struct {
	TemplateType       string          `json:"overlay_template_type"`
	MessageOverlay     OverlayTemplate `json:"message_overlay"`
	DescriptionOverlay OverlayTemplate `json:"description_overlay"`
}

type OverlayTemplate struct {
	Template string `json:"template"`
}

type CustomContentTemplateOverlay struct {
	ServiceID            string                `json:"service_id" tf:"service_id"`
	AlertSourceShortname string                `json:"alert_source_shortname" tf:"alert_source_shortname"`
	OverlayTemplateType  string                `json:"overlay_template_type" tf:"-"`
	Overlay              CustomContentTemplate `json:"overlay" tf:"-"`
}

type CustomContentTemplate struct {
	Message     string `json:"message" tf:"-"`
	Description string `json:"description" tf:"-"`
}

func (w *CustomContentTemplateOverlay) Encode() (tf.M, error) {
	m, err := tf.Encode(w)
	if err != nil {
		return nil, err
	}

	m["message_template"] = w.Overlay.Message
	m["description_template"] = w.Overlay.Description

	return m, nil
}

func (client *Client) CreateOrUpdateCustomContentTemplateOverlay(ctx context.Context, serviceID, alertSourceShortname string, req *CustomContentTemplateOverlayReq) (*CustomContentTemplateOverlay, error) {
	url := fmt.Sprintf("%s/services/%s/overlays/custom-content/%s", client.BaseURLV3, serviceID, alertSourceShortname)
	return Request[CustomContentTemplateOverlayReq, CustomContentTemplateOverlay](http.MethodPut, url, client, ctx, req)
}

func (client *Client) GetCustomContentTemplateOverlay(ctx context.Context, serviceID, alertSourceShortname string) (*CustomContentTemplateOverlay, error) {
	url := fmt.Sprintf("%s/services/%s/overlays/custom-content/%s", client.BaseURLV3, serviceID, alertSourceShortname)
	return Request[any, CustomContentTemplateOverlay](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) DeleteCustomContentTemplateOverlay(ctx context.Context, serviceID, alertSourceShortname string) error {
	url := fmt.Sprintf("%s/services/%s/overlays/custom-content/%s", client.BaseURLV3, serviceID, alertSourceShortname)
	_, err := Request[any, any](http.MethodDelete, url, client, ctx, nil)
	return err
}
