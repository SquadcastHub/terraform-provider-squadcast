package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func resourceWorkflowAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowActionCreate,
		ReadContext:   resourceWorkflowActionRead,
		UpdateContext: resourceWorkflowActionUpdate,
		DeleteContext: resourceWorkflowActionDelete,
		Schema: map[string]*schema.Schema{
			"workflow_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"sq_add_incident_note"}, false),
			},
			"note": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceWorkflowActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	tflog.Info(ctx, "Creating a new workflow action", tf.M{
		"name": d.Get("name").(string),
	})

	client := meta.(*api.Client)

	workflowAction := &api.WorkflowAction{
		Name: d.Get("name").(string),
		Data: api.WorkflowActionData{
			Note: d.Get("note").(string),
		},
	}

	workflowID := d.Get("workflow_id").(string)

	workflowActionResponse, err := client.CreateWorkflowAction(ctx, workflowID, workflowAction)
	if err != nil {
		return diag.FromErr(err)
	}

	workflowActionID := strconv.FormatUint(uint64(workflowActionResponse.ID), 10)
	d.SetId(workflowActionID)

	return resourceWorkflowActionRead(ctx, d, meta)
}

func resourceWorkflowActionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Reading workflow action", tf.M{
		"action_id":   d.Id(),
		"workflow_id": d.Get("workflow_id").(string),
	})

	workflowID := d.Get("workflow_id").(string)
	workflowActionID := d.Id()

	workflowAction, err := client.GetWorkflowActionById(ctx, workflowID, workflowActionID)
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO: Check if this is the correct way to encode the data
	workFlowActionEncoded := &api.WorkflowAction{
		Name:       workflowAction.Name,
		WorkflowID: workflowID,
		Data: api.WorkflowActionData{
			Note: workflowAction.Data.Note,
		},
	}

	if err = tf.EncodeAndSet(workFlowActionEncoded, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWorkflowActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*api.Client)

	tflog.Info(ctx, "Updating workflow action", tf.M{
		"action_id": d.Id(),
	})

	workflowAction := &api.WorkflowAction{
		WorkflowID: d.Get("workflow_id").(string),
		Name:       d.Get("name").(string),
		Data: api.WorkflowActionData{
			Note: d.Get("note").(string),
		},
	}

	workflowID := d.Get("workflow_id").(string)

	_, err := client.UpdateWorkflowAction(ctx, workflowID, d.Id(), workflowAction)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceWorkflowActionRead(ctx, d, meta)
}

func resourceWorkflowActionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Deleting workflow action", tf.M{
		"id": d.Id(),
	})

	workflowID := d.Get("workflow_id").(string)

	_, err := client.DeleteWorkflowAction(ctx, workflowID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
