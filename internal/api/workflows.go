package api

import (
	"context"
	"fmt"
	"net/http"
)

type Workflow struct {
	ID          string      `json:"ID" tf:"id"`
	Title       string      `json:"title" tf:"title"`
	Description string      `json:"description" tf:"description"`
	OwnerID     string      `json:"ownerID" tf:"owner_id"`
	Enabled     bool        `json:"enabled" tf:"enabled"`
	Trigger     string      `json:"trigger" tf:"trigger"`
	Filters     []Filters   `json:"filters" tf:"filters"`
	EntityOwner EntityOwner `json:"entityOwner" tf:"entity_owner"`
	Tags        []*Tag      `json:"tags" tf:"tags"`
}

type Filters struct {
	Fields Field  `json:"fields" tf:"fields"`
	Type   string `json:"type" tf:"type"`
}

type Field struct {
	Value string `json:"value" tf:"value"`
}

func (client *Client) CreateWorkflow(ctx context.Context, workflowReq *Workflow) (*Workflow, error) {
	url := fmt.Sprintf("%s/workflows", client.BaseURLV3)
	return Request[Workflow, Workflow](http.MethodPost, url, client, ctx, workflowReq)
}
