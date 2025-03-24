package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type DedupKeyOverlayReq struct {
	TemplateType    string   `json:"overlay_template_type"`
	DedupKeyOverlay DedupKey `json:"dedup_key_overlay"`
}

type DedupKeyOverlay struct {
	ServiceID            string   `json:"service_id" tf:"service_id"`
	AlertSourceShortname string   `json:"alert_source_shortname" tf:"alert_source_shortname"`
	OverlayTemplateType  string   `json:"overlay_template_type" tf:"-"`
	Overlay              DedupKey `json:"overlay" tf:"-"`
}

type DedupKey struct {
	Template       string `json:"template" tf:"-"`
	DurationInMins uint   `json:"duration" tf:"-"`
}

func (w *DedupKeyOverlay) Encode() (tf.M, error) {
	m, err := tf.Encode(w)
	if err != nil {
		return nil, err
	}

	m["dedup_key_overlay_template"] = w.Overlay.Template
	m["duration"] = w.Overlay.DurationInMins

	return m, nil
}

func (client *Client) CreateOrUpdateDedupKeyOverlay(ctx context.Context, serviceID, alertSourceShortname string, req *DedupKeyOverlayReq) (*DedupKeyOverlay, error) {
	url := fmt.Sprintf("%s/services/%s/overlays/dedup-key/%s", client.BaseURLV3, serviceID, alertSourceShortname)
	return Request[DedupKeyOverlayReq, DedupKeyOverlay](http.MethodPut, url, client, ctx, req)
}

func (client *Client) GetDedupKeyOverlay(ctx context.Context, serviceID, alertSourceShortname string) (*DedupKeyOverlay, error) {
	url := fmt.Sprintf("%s/services/%s/overlays/dedup-key/%s", client.BaseURLV3, serviceID, alertSourceShortname)
	return Request[any, DedupKeyOverlay](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) DeleteDedupKeyOverlay(ctx context.Context, serviceID, alertSourceShortname string) error {
	url := fmt.Sprintf("%s/services/%s/overlays/dedup-key/%s", client.BaseURLV3, serviceID, alertSourceShortname)
	_, err := Request[any, any](http.MethodDelete, url, client, ctx, nil)
	return err
}
