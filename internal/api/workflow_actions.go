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
	Note string `json:"note" tf:"note"`
	// Runbooks
	Runbooks []string `json:"runbooks" tf:"runbooks"`
	// SLO
	SLO  int      `json:"slo" tf:"slo"`
	SLIs []string `json:"slis" tf:"slis"`
	// Communication channels
	Channels []Channels `json:"channels" tf:"channels"`
	// Incident priority
	Priority string `json:"priority" tf:"priority"`
	// Http request
	Method  string    `json:"method" tf:"method"`
	URL     string    `json:"url" tf:"url"`
	Body    string    `json:"body" tf:"body"`
	Headers []Headers `json:"headers" tf:"headers"`
	// Send Email
	To      []string `json:"to" tf:"to"`
	Subject string   `json:"subject" tf:"subject"`
	// body is needed for email as well
	// Trigger Manual Webhook
	WebhookID string `json:"id" tf:"webhook_id"`
	// Status page
	StatusPageID       int                  `json:"status_page_id" tf:"status_page_id"`
	IssueTitle         string               `json:"issue_title" tf:"issue_title"`
	PageStatusID       int                  `json:"page_status_id" tf:"page_status_id"`
	ComponentAndImpact []ComponentAndImpact `json:"component_and_impact" tf:"component_and_impact"`
	StatusAndMessage   []StatusAndMessage   `json:"status_and_message" tf:"status_and_message"`
}

type ComponentAndImpact struct {
	ComponentID    int `json:"component_id" tf:"component_id"`
	ImpactStatusID int `json:"impact_status_id" tf:"impact_status_id"`
}

type StatusAndMessage struct {
	StatusID int      `json:"status_id" tf:"status_id"`
	Messages []string `json:"messages" tf:"messages"`
}

type Channels struct {
	ChannelType string `json:"type" tf:"type"`
	Link        string `json:"link" tf:"link"`
	DisplayText string `json:"display_text" tf:"display_text"`
}

type Headers struct {
	Key   string `json:"key" tf:"key"`
	Value string `json:"value" tf:"value"`
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
	Method   string     `json:"method" tf:"method"`
	URL      string     `json:"url" tf:"url"`
	Body     string     `json:"body" tf:"body"`
	Headers  []struct {
		Key   string `json:"key" tf:"key"`
		Value string `json:"value" tf:"value"`
	} `json:"headers" tf:"headers"`
	To                 []string                `json:"to" tf:"to"`
	Subject            string                  `json:"subject" tf:"subject"`
	WebhookID          string                  `json:"id" tf:"webhook_id"`
	StatusPageID       int                     `json:"status_page_id" tf:"status_page_id"`
	IssueTitle         string                  `json:"issue_title" tf:"issue_title"`
	PageStatusID       int                     `json:"page_status_id" tf:"page_status_id"`
	ComponentAndImpact []ComponentAndImpactRes `json:"component_and_impact" tf:"component_and_impact"`
	StatusAndMessage   []StatusAndMessage      `json:"status_and_message" tf:"status_and_message"`
}

type ComponentAndImpactRes struct {
	ComponentID    int `json:"component_id" tf:"component_id"`
	ImpactStatusID int `json:"impact_status_id" tf:"impact_status_id"`
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
