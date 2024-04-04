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
				Type:        schema.TypeString,
				Description: "The ID of the workflow to which this action belongs",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the action",
				Required:    true,
				ValidateFunc: validation.StringInSlice([]string{"sq_add_incident_note", "sq_attach_runbooks",
					"sq_mark_incident_slo_affecting", "sq_add_communication_channel", "sq_update_incident_priority"}, false),
			},
			"note": {
				Type:        schema.TypeString,
				Description: "The note to be added to the incident",
				Optional:    true,
			},
			"runbooks": {
				Type:        schema.TypeList,
				Description: "The runbooks to be added to the incident",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"slo": {
				Type:        schema.TypeInt,
				Description: "ID of the SLO to be added to the incident",
				Optional:    true,
			},
			"slis": {
				Type:        schema.TypeList,
				Description: "The SLIs to be added to the incident",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"channels": {
				Type:        schema.TypeList,
				Description: "The communication channels to be added to the incident",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Description:  "The type of the communication channel",
							ValidateFunc: validation.StringInSlice([]string{"chat_room", "video_conference", "other"}, false),
							Required:     true,
						},
						"link": {
							Type:        schema.TypeString,
							Description: "The link of the communication channel",
							Required:    true,
						},
						"display_text": {
							Type:        schema.TypeString,
							Description: "The display text of the communication channel",
							Required:    true,
						},
					},
				},
			},
			"priority": {
				Type:         schema.TypeString,
				Description:  "The priority of the incident",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"P1", "P2", "P3", "P4", "P5", "UNSET"}, false),
			},
		},
	}
}

func resourceWorkflowActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*api.Client)

	tflog.Info(ctx, "Creating a new workflow action", tf.M{
		"name":        d.Get("name").(string),
		"worfklow_id": d.Get("workflow_id").(string),
	})

	runbooks := tf.ListToSlice[string](d.Get("runbooks"))
	channels := make([]api.Channels, 0)

	if err := Decode(d.Get("channels"), &channels); err != nil {
		return diag.FromErr(err)
	}

	workflowAction := &api.WorkflowAction{
		Name: d.Get("name").(string),
		Data: api.WorkflowActionData{
			Note:     d.Get("note").(string),
			SLO:      d.Get("slo").(int),
			SLIs:     tf.ListToSlice[string](d.Get("slis")),
			Priority: d.Get("priority").(string),
			Runbooks: runbooks,
			Channels: channels,
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

func resourceWorkflowActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*api.Client)

	tflog.Info(ctx, "Updating workflow action", tf.M{
		"worfklow_id": d.Get("workflow_id").(string),
		"action_id":   d.Id(),
	})

	runbooks := tf.ListToSlice[string](d.Get("runbooks"))
	channels := make([]api.Channels, 0)

	if err := Decode(d.Get("channels"), &channels); err != nil {
		return diag.FromErr(err)
	}

	workflowAction := &api.WorkflowAction{
		Name: d.Get("name").(string),
		Data: api.WorkflowActionData{
			Note:     d.Get("note").(string),
			SLO:      d.Get("slo").(int),
			SLIs:     tf.ListToSlice[string](d.Get("slis")),
			Priority: d.Get("priority").(string),
			Runbooks: runbooks,
			Channels: channels,
		},
	}

	workflowID := d.Get("workflow_id").(string)

	_, err := client.UpdateWorkflowAction(ctx, workflowID, d.Id(), workflowAction)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceWorkflowActionRead(ctx, d, meta)
}

func resourceWorkflowActionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Reading workflow action", tf.M{
		"name":        d.Get("name").(string),
		"action_id":   d.Id(),
		"worfklow_id": d.Get("workflow_id").(string),
	})

	workflowID := d.Get("workflow_id").(string)
	workflowActionID := d.Id()

	workflowAction, err := client.GetWorkflowActionById(ctx, workflowID, workflowActionID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(workflowAction, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWorkflowActionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Deleting workflow action", tf.M{
		"worfklow_id": d.Get("workflow_id").(string),
		"action_id":   d.Id(),
	})

	workflowID := d.Get("workflow_id").(string)

	_, err := client.DeleteWorkflowAction(ctx, workflowID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
