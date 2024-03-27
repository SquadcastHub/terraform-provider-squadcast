package api

import (
	"context"
	"fmt"
	"net/http"
)

type Workflow struct {
	ID          uint        `json:"id" tf:"id,omitempty"`
	Title       string      `json:"title" tf:"title"`
	Description string      `json:"description" tf:"description"`
	OwnerID     string      `json:"owner_id" tf:"owner_id"`
	Enabled     bool        `json:"enabled" tf:"enabled"`
	Trigger     string      `json:"trigger" tf:"trigger"`
	Filters     []Filters   `json:"filters" tf:"filters"`
	EntityOwner EntityOwner `json:"entity_owner" tf:"entity_owner"`
	Tags        []Tag       `json:"tags,omitempty" tf:"tags"`
}

type Filters struct {
	Fields Field  `json:"fields" tf:"fields"`
	Type   string `json:"type" tf:"type"`
}

type Field struct {
	Value string `json:"value" tf:"value"`
}

func (client *Client) CreateWorkflow(ctx context.Context, workflowReq *Workflow) (*Workflow, error) {
	// url := "https://webhook.site/92de1537-f934-4c4b-b521-0a8ca6d0638f"
	url := fmt.Sprintf("%s/workflows", client.BaseURLV3)
	return Request[Workflow, Workflow](http.MethodPost, url, client, ctx, workflowReq)
}
