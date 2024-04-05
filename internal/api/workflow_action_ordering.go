package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type WorkflowActionOrdering struct {
	WorkflowID  string `json:"workflow_id" tf:"workflow_id"`
	ActionOrder []int  `json:"action_order" tf:"action_order"`
}

func (wao *WorkflowActionOrdering) Encode() (tf.M, error) {
	return tf.Encode(wao)
}

func (client *Client) UpdateWorkflowActionOrdering(ctx context.Context, workflowID string, workflowActionOrder *WorkflowActionOrdering) (*WorkflowActionOrdering, error) {
	url := fmt.Sprintf("%s/workflows/%s/actions/reorder", client.BaseURLV3, workflowID)
	return Request[WorkflowActionOrdering, WorkflowActionOrdering](http.MethodPatch, url, client, ctx, workflowActionOrder)
}

func (client *Client) GetWorkflowActionOrdering(ctx context.Context, workflowID string) (*WorkflowActionOrdering, error) {
	url := fmt.Sprintf("%s/workflows/%s", client.BaseURLV3, workflowID)
	workflow, err := Request[any, Workflow](http.MethodGet, url, client, ctx, nil)
	if err != nil {
		return nil, err
	}

	order := make([]int, 0)
	for _, action := range workflow.Actions {
		if action.ID == 0 {
			return nil, fmt.Errorf("action order is not set for workflow %s", workflowID)
		}
		order = append(order, action.ID)
	}

	return &WorkflowActionOrdering{
		WorkflowID:  workflowID,
		ActionOrder: order,
	}, nil
}
