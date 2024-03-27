package api

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type Workflow struct {
	ID          uint   `json:"id" tf:"id,omitempty"`
	Title       string `json:"title" tf:"title"`
	Description string `json:"description" tf:"description"`
	OwnerID     string `json:"owner_id" tf:"owner_id"`
	Enabled     bool   `json:"enabled" tf:"enabled"`
	Trigger     string `json:"trigger" tf:"trigger"`
	// Filters     []*Filters  `json:"filters,omitempty" tf:"filters"`
	EntityOwner EntityOwner             `json:"entity_owner" tf:"entity_owner"`
	Tags        map[string]TagWithColor `json:"tags" tf:"-"`
}

// TODO: name it better
type TagWithColor struct {
	Value string `json:"value" tf:"value"`
	Color string `json:"color" tf:"color"`
	Key   string `json:"key" tf:"key"`
}

type Filters struct {
	Fields Field  `json:"fields" tf:"fields"`
	Type   string `json:"type" tf:"type"`
}

type Field struct {
	Value string `json:"value" tf:"value"`
}

func (twc *TagWithColor) Encode() (tf.M, error) {
	return tf.Encode(twc)
}

func (w *Workflow) Encode() (tf.M, error) {
	m, err := tf.Encode(w)
	if err != nil {
		return nil, err
	}

	tags := make([]any, 0, len(w.Tags))

	keys := make([]string, 0, len(w.Tags))
	for k := range w.Tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := w.Tags[k]
		mtag, err := v.Encode()
		if err != nil {
			return nil, err
		}
		mtag["key"] = k
		tags = append(tags, mtag)
	}
	m["tags"] = tags

	m["entity_owner"] = tf.List(tf.M{
		"id":   w.EntityOwner.ID,
		"type": w.EntityOwner.Type,
	})

	return m, nil
}

func (client *Client) CreateWorkflow(ctx context.Context, workflowReq *Workflow) (*Workflow, error) {
	// url := "https://webhook.site/92de1537-f934-4c4b-b521-0a8ca6d0638f"
	url := fmt.Sprintf("%s/workflows", client.BaseURLV3)
	return Request[Workflow, Workflow](http.MethodPost, url, client, ctx, workflowReq)
}

func (client *Client) GetWorkflowById(ctx context.Context, id string) (*Workflow, error) {
	url := fmt.Sprintf("%s/workflows/%s", client.BaseURLV3, id)
	return Request[any, Workflow](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) UpdateWorkflow(ctx context.Context, id string, workflowReq *Workflow) (*Workflow, error) {
	url := fmt.Sprintf("%s/workflows/%s", client.BaseURLV3, id)
	return Request[Workflow, Workflow](http.MethodPut, url, client, ctx, workflowReq)
}

func (client *Client) DeleteWorkflow(ctx context.Context, id string) (*any, error) {
	url := fmt.Sprintf("%s/workflows/%s", client.BaseURLV3, id)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
