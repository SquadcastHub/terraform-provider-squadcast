package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type Workflow struct {
	ID          uint   `json:"id,omitempty" tf:"id,omitempty"`
	Title       string `json:"title" tf:"title"`
	Description string `json:"description" tf:"description"`
	OwnerID     string `json:"owner_id" tf:"owner_id"`
	Enabled     bool   `json:"enabled" tf:"enabled"`
	Trigger     string `json:"trigger" tf:"trigger"`
	// Filters     []*Filters  `json:"filters,omitempty" tf:"filters"`
	EntityOwner EntityOwner    `json:"entity_owner" tf:"entity_owner"`
	Tags        []*WorkflowTag `json:"tags,omitempty" tf:"-"`

	// Should be used only for action ordering resource
	// Hence we are not encoding this field (tf:"-")
	Actions []*WorkflowAction `json:"actions,omitempty" tf:"-"`
}

type WorkflowTag struct {
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

func (twc *WorkflowTag) Encode() (tf.M, error) {
	return tf.Encode(twc)
}

func (w *Workflow) Encode() (tf.M, error) {
	m, err := tf.Encode(w)
	if err != nil {
		return nil, err
	}

	tagsEncoded, terr := tf.EncodeSlice(w.Tags)
	if terr != nil {
		return nil, terr
	}
	m["tags"] = tagsEncoded

	m["entity_owner"] = tf.List(tf.M{
		"id":   w.EntityOwner.ID,
		"type": w.EntityOwner.Type,
	})

	return m, nil
}

func (client *Client) CreateWorkflow(ctx context.Context, workflowReq *Workflow) (*Workflow, error) {
	url := fmt.Sprintf("%s/workflows", client.BaseURLV3)
	return Request[Workflow, Workflow](http.MethodPost, url, client, ctx, workflowReq)
}

func (client *Client) GetWorkflowById(ctx context.Context, id string) (*Workflow, error) {
	url := fmt.Sprintf("%s/workflows/%s", client.BaseURLV3, id)
	return Request[any, Workflow](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) UpdateWorkflow(ctx context.Context, id string, workflowReq *Workflow) (*Workflow, error) {
	url := fmt.Sprintf("%s/workflows/%s", client.BaseURLV3, id)
	return Request[Workflow, Workflow](http.MethodPatch, url, client, ctx, workflowReq)
}

func (client *Client) DeleteWorkflow(ctx context.Context, id string) (*any, error) {
	url := fmt.Sprintf("%s/workflows/%s", client.BaseURLV3, id)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}

func Decode(input any, output any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:               output,
		TagName:              tf.EncoderStructTag,
		ZeroFields:           true,
		IgnoreUntaggedFields: true,
	})
	if err != nil {
		return err
	}

	err = decoder.Decode(input)
	if err != nil {
		return err
	}

	return nil
}
