package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type WorkflowAction struct {
	Name string             `json:"name" tf:"name"`
	Data WorkflowActionData `json:"data" tf:"data"`
}

type WorkflowActionData struct {
	Note string `json:"note" tf:"note"`
}

type WorkflowActionResponse struct {
	Data struct {
		Data struct {
			Note string `json:"note"`
		} `json:"data"`
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
}

func (twc *WorkflowActionData) Encode() (tf.M, error) {
	return tf.Encode(twc)
}

func (w *WorkflowAction) Encode() (tf.M, error) {
	return tf.Encode(w)
}

func (client *Client) CreateWorkflowAction(ctx context.Context, workflowID string, workflowAction *WorkflowAction) (*WorkflowActionResponse, error) {
	url := fmt.Sprintf("%s/workflows/%s/actions", client.BaseURLV3, workflowID)
	return Request[WorkflowAction, WorkflowActionResponse](http.MethodPost, url, client, ctx, workflowAction)
}

func (client *Client) GetWorkflowActionById(ctx context.Context, workflowID, actionID string) (*WorkflowActionResponse, error) {
	url := fmt.Sprintf("%s/workflows/%s/actions/%s", client.BaseURLV3, workflowID, actionID)
	return Request[WorkflowAction, WorkflowActionResponse](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) UpdateWorkflowAction(ctx context.Context, workflowID, actionID string, workflowAction *WorkflowAction) (*WorkflowActionResponse, error) {
	url := fmt.Sprintf("%s/workflows/%s/actions/%s", client.BaseURLV3, workflowID, actionID)
	return Request[WorkflowAction, WorkflowActionResponse](http.MethodPut, url, client, ctx, workflowAction)
}

func (client *Client) DeleteWorkflowAction(ctx context.Context, workflowID, actionID string) (*any, error) {
	url := fmt.Sprintf("%s/workflows/%s/actions/%s", client.BaseURLV3, workflowID, actionID)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
