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

	// Filters     *Filters       `json:"filters,omitempty" tf:"filters"`
	Filters *HighLevelFilter `json:"filters,omitempty" tf:"filters"`

	EntityOwner EntityOwner    `json:"entity_owner" tf:"entity_owner"`
	Tags        []*WorkflowTag `json:"tags,omitempty" tf:"-"`

	// Should be used only for action ordering resource
	// Hence we are not encoding this field (tf:"-")
	Actions []*WorkflowAction `json:"actions,omitempty" tf:"-"`
}

type HighLevelFilter struct {
	Condition string        `json:"condition" tf:"condition"`
	Filters   []FilterGroup `json:"filters" tf:"filters"`
}

type FilterGroup struct {
	Condition string          `json:"condition,omitempty" tf:"condition"`
	Type      string          `json:"type,omitempty" tf:"type"`
	Value     string          `json:"value,omitempty" tf:"value"`
	Filters   []*ChildFilters `json:"filters,omitempty" tf:"filters"`
}

type ChildFilters struct {
	Type  string `json:"type" tf:"type"`
	Key   string `json:"key" tf:"key"`
	Value string `json:"value" tf:"value"`
}

type WorkflowTag struct {
	Value string `json:"value" tf:"value"`
	Color string `json:"color" tf:"color"`
	Key   string `json:"key" tf:"key"`
}

func (f FilterGroup) Encode() (tf.M, error) {
	return tf.Encode(f)
}

func (f *ChildFilters) Encode() (tf.M, error) {
	return tf.Encode(f)
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

	if w.Filters != nil {

		filters := tf.List(tf.M{})
		for _, filter := range w.Filters.Filters {

			fData := tf.M{}
			if filter.Condition != "" {
				fData["condition"] = filter.Condition
			}

			if filter.Type != "" {
				fData["type"] = filter.Type
			}

			if filter.Value != "" {
				fData["value"] = filter.Value
			}

			if filter.Filters != nil {
				childFilters := tf.List(tf.M{})
				for _, childFilter := range filter.Filters {
					childFilterData := tf.M{}
					if childFilter.Type != "" {
						childFilterData["type"] = childFilter.Type
					}
					if childFilter.Key != "" {
						childFilterData["key"] = childFilter.Key
					}
					if childFilter.Value != "" {
						childFilterData["value"] = childFilter.Value
					}
					childFilters = append(childFilters, childFilterData)
				}
				fData["filters"] = childFilters
			}

			filters = append(filters, fData)
		}

		m["filters"] = tf.List(tf.M{
			"condition": w.Filters.Condition,
			"filters":   filters,
		})

		// tf.List(tf.M{
		// 	"condition": "",
		// 	"filters": tf.List(tf.M{
		// 		"condition": "",
		// 		"type":      "",
		// 		"value":     "",
		// 		"filters": tf.List(tf.M{
		// 			"type":  "",
		// 			"key":   "",
		// 			"value": "",
		// 		}),
		// 	}),
		// })
	}

	m["entity_owner"] = tf.List(tf.M{
		"id":   w.EntityOwner.ID,
		"type": w.EntityOwner.Type,
	})

	return m, nil
}

func (client *Client) CreateWorkflow(ctx context.Context, workflowReq *Workflow) (*Workflow, error) {
	url := fmt.Sprintf("%s/workflows", client.BaseURLV3)
	// url := "https://webhook.site/9cc3df1e-7d1f-458d-9635-305300a9c1a5"
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
