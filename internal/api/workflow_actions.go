package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

// WorkflowAction represents a workflow action request & terraform resource
type WorkflowAction struct {
	ID         int                `json:"id" tf:"id"`
	Name       string             `json:"name" tf:"name"`
	WorkflowID int                `json:"workflow_id" tf:"workflow_id"`
	Data       WorkflowActionData `json:"data" tf:"-"`
}

type WorkflowActionData struct {
	Note     string     `json:"note" tf:"note"`
	Runbooks []string   `json:"runbooks" tf:"runbooks"`
	SLO      int        `json:"slo" tf:"slo"`
	SLIs     []string   `json:"slis" tf:"slis"`
	Channels []Channels `json:"channels" tf:"channels"`
	Priority string     `json:"priority" tf:"priority"`
}

type Channels struct {
	ChannelType string `json:"type" tf:"type"`
	Link        string `json:"link" tf:"link"`
	DisplayText string `json:"display_text" tf:"display_text"`
}

// WorkflowActionRes represents a workflow action response
type WorkflowActionRes struct {
	ID         int                   `json:"id" tf:"id"`
	Name       string                `json:"name" tf:"name"`
	WorkflowID int                   `json:"workflow_id" tf:"workflow_id"`
	Data       WorkflowActionDataRes `json:"data" tf:"-"`
}

type WorkflowActionDataRes struct {
	Note     string   `json:"note" tf:"note"`
	SLO      string   `json:"name" tf:"slo"` // response sends id where request sends name
	SLIs     []string `json:"slis" tf:"slis"`
	Runbooks []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"runbooks"`
	Channels []Channels `json:"channels" tf:"channels"`
	Priority string     `json:"priority" tf:"priority"`
}

// TODO: Check if we even need this??
func (twc *WorkflowActionDataRes) Encode() (tf.M, error) {
	return tf.Encode(twc)
}
func (w *WorkflowActionRes) Encode() (tf.M, error) {
	m, err := tf.Encode(w)
	if err != nil {
		return nil, err
	}

	m["workflow_id"] = fmt.Sprintf("%d", w.WorkflowID)
	return m, nil
}

func (client *Client) CreateWorkflowAction(ctx context.Context, workflowID string, workflowAction *WorkflowAction) (*WorkflowActionRes, error) {
	// url := "https://webhook.site/3ee9072d-5587-4bb3-95af-bed7db01fffa"
	url := fmt.Sprintf("%s/workflows/%s/actions", client.BaseURLV3, workflowID)
	return Request[WorkflowAction, WorkflowActionRes](http.MethodPost, url, client, ctx, workflowAction)
}

func (client *Client) GetWorkflowActionById(ctx context.Context, workflowID, actionID string) (*WorkflowActionRes, error) {
	url := fmt.Sprintf("%s/workflows/%s/actions/%s", client.BaseURLV3, workflowID, actionID)
	return Request[WorkflowAction, WorkflowActionRes](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) UpdateWorkflowAction(ctx context.Context, workflowID, actionID string, workflowAction *WorkflowAction) (*WorkflowActionRes, error) {
	url := fmt.Sprintf("%s/workflows/%s/actions/%s", client.BaseURLV3, workflowID, actionID)
	return Request[WorkflowAction, WorkflowActionRes](http.MethodPatch, url, client, ctx, workflowAction)
}

func (client *Client) DeleteWorkflowAction(ctx context.Context, workflowID, actionID string) (*any, error) {
	url := fmt.Sprintf("%s/workflows/%s/actions/%s", client.BaseURLV3, workflowID, actionID)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
