package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func resourceWorkflowActionOrdering() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowActionOrderingUpdate,
		ReadContext:   resourceWorkflowActionOrderingRead,
		UpdateContext: resourceWorkflowActionOrderingUpdate,
		DeleteContext: resourceWorkflowActionOrderingDelete,
		// Import: Not required for this resource
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of this resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workflow_id": {
				Type:        schema.TypeString,
				Description: "The ID of the workflow",
				Required:    true,
			},
			"action_order": {
				Type:        schema.TypeList,
				Description: "The order of actions in the workflow",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceWorkflowActionOrderingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Creating workflow action ordering", tf.M{
		"workflow_id": d.Get("workflow_id").(string),
	})

	action_order := d.Get("action_order").([]interface{})
	action_ids := make([]int, len(action_order))
	for i, v := range action_order {
		action_ids[i] = v.(int)
	}

	workflowActionOrdering := &api.WorkflowActionOrdering{
		WorkflowID:  d.Get("workflow_id").(string),
		ActionOrder: action_ids,
	}

	_, err := client.UpdateWorkflowActionOrdering(ctx, workflowActionOrdering.WorkflowID, workflowActionOrdering)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(workflowActionOrdering.WorkflowID)

	return resourceWorkflowActionOrderingRead(ctx, d, meta)
}

func resourceWorkflowActionOrderingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Reading workflow action ordering", tf.M{
		"workflow_id": d.Id(),
	})

	workflowActionOrdering, err := client.GetWorkflowActionOrdering(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(workflowActionOrdering, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWorkflowActionOrderingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
