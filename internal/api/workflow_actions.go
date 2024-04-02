package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type WorkflowAction struct {
	ID         int                `json:"id" tf:"id"`
	Name       string             `json:"name" tf:"name"`
	WorkflowID int                `json:"workflow_id" tf:"workflow_id"`
	Data       WorkflowActionData `json:"data" tf:"-"`
}

type WorkflowActionData struct {
	Note string `json:"note" tf:"note"`
}

// TODO: Check if we even need this??
func (twc *WorkflowActionData) Encode() (tf.M, error) {
	return tf.Encode(twc)
}

func (w *WorkflowAction) Encode() (tf.M, error) {
	m, err := tf.Encode(w)
	if err != nil {
		return nil, err
	}

	m["workflow_id"] = fmt.Sprintf("%d", w.WorkflowID)
	return m, nil
}

func (client *Client) CreateWorkflowAction(ctx context.Context, workflowID string, workflowAction *WorkflowAction) (*WorkflowAction, error) {
	url := fmt.Sprintf("%s/workflows/%s/actions", client.BaseURLV3, workflowID)
	return Request[WorkflowAction, WorkflowAction](http.MethodPost, url, client, ctx, workflowAction)
}

func (client *Client) GetWorkflowActionById(ctx context.Context, workflowID, actionID string) (*WorkflowAction, error) {
	url := fmt.Sprintf("%s/workflows/%s/actions/%s", client.BaseURLV3, workflowID, actionID)
	return Request[WorkflowAction, WorkflowAction](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) UpdateWorkflowAction(ctx context.Context, workflowID, actionID string, workflowAction *WorkflowAction) (*WorkflowAction, error) {
	url := fmt.Sprintf("%s/workflows/%s/actions/%s", client.BaseURLV3, workflowID, actionID)
	return Request[WorkflowAction, WorkflowAction](http.MethodPatch, url, client, ctx, workflowAction)
}

func (client *Client) DeleteWorkflowAction(ctx context.Context, workflowID, actionID string) (*any, error) {
	url := fmt.Sprintf("%s/workflows/%s/actions/%s", client.BaseURLV3, workflowID, actionID)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
