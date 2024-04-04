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
	Note     string   `json:"note" tf:"note"`
	Runbooks []string `json:"runbooks" tf:"runbooks"`
}

// Runbook action
type WorkflowActionRunbookRes struct {
	ID         int                            `json:"id" tf:"id"`
	Name       string                         `json:"name" tf:"name"`
	WorkflowID int                            `json:"workflow_id" tf:"workflow_id"`
	Data       []WorkflowActionRunbookDataRes `json:"data" tf:"-"`
}

type WorkflowActionRunbookDataRes struct {
	ID   string `json:"id" tf:"id"`
	Name string `json:"name" tf:"name"`
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

// Runbook action response
func (w *WorkflowActionRunbookRes) Encode() (tf.M, error) {
	m, err := tf.Encode(w)
	if err != nil {
		return nil, err
	}

	m["workflow_id"] = fmt.Sprintf("%d", w.WorkflowID)
	return m, nil
}

func (twc *WorkflowActionRunbookDataRes) Encode() (tf.M, error) {
	return tf.Encode(twc)
}

func (client *Client) CreateWorkflowAction(ctx context.Context, workflowID string, workflowAction *WorkflowAction) (*WorkflowAction, error) {
	url := fmt.Sprintf("%s/workflows/%s/actions", client.BaseURLV3, workflowID)
	// url := "https://webhook.site/3ee9072d-5587-4bb3-95af-bed7db01fffa"
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

// Runbook action
func (client *Client) CreateRunbookWorkflowAction(ctx context.Context, workflowID string, workflowAction *WorkflowAction) (*WorkflowActionRunbookRes, error) {
	url := fmt.Sprintf("%s/workflows/%s/actions", client.BaseURLV3, workflowID)
	// url := "https://webhook.site/3ee9072d-5587-4bb3-95af-bed7db01fffa"
	return Request[WorkflowAction, WorkflowActionRunbookRes](http.MethodPost, url, client, ctx, workflowAction)
}

func (client *Client) GetRunbookWorkflowActionById(ctx context.Context, workflowID, actionID string) (*WorkflowActionRunbookRes, error) {
	url := fmt.Sprintf("%s/workflows/%s/actions/%s", client.BaseURLV3, workflowID, actionID)
	return Request[WorkflowAction, WorkflowActionRunbookRes](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) UpdateRunbookWorkflowAction(ctx context.Context, workflowID, actionID string, workflowAction *WorkflowAction) (*WorkflowActionRunbookRes, error) {
	url := fmt.Sprintf("%s/workflows/%s/actions/%s", client.BaseURLV3, workflowID, actionID)
	return Request[WorkflowAction, WorkflowActionRunbookRes](http.MethodPatch, url, client, ctx, workflowAction)
}

func (client *Client) DeleteWorkflowAction(ctx context.Context, workflowID, actionID string) (*any, error) {
	url := fmt.Sprintf("%s/workflows/%s/actions/%s", client.BaseURLV3, workflowID, actionID)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
