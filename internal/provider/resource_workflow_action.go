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
					"sq_mark_incident_slo_affecting", "sq_add_communication_channel", "sq_update_incident_priority",
					"sq_make_http_call", "sq_send_email", "sq_trigger_manual_webhook", "sq_add_status_page_issue", "jira_create_ticket",
					"slack_create_incident_channel", "slack_archive_channel", "slack_message_channel"}, false),
			},
			// Add Notes Action
			"note": {
				Type:        schema.TypeString,
				Description: "The note to be added to the incident",
				Optional:    true,
			},
			// Attach Runbooks Action
			"runbooks": {
				Type:        schema.TypeList,
				Description: "The runbooks to be added to the incident",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// SLO Action
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
			// Communication Channel Action
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
			// Incident Priority Action
			"priority": {
				Type:         schema.TypeString,
				Description:  "The priority of the incident",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"P1", "P2", "P3", "P4", "P5", "UNSET"}, false),
			},
			// HTTP Call Action
			"method": {
				Type:         schema.TypeString,
				Description:  "The HTTP method to be used for the call",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"GET", "POST", "PUT", "PATCH", "DELETE"}, false),
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The URL to be called",
				Optional:    true,
			},
			"headers": {
				Type:        schema.TypeList,
				Description: "The headers to be sent with the request",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Description: "The key of the header",
							Required:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "The value of the header",
							Required:    true,
						},
					},
				},
			},
			"body": {
				Type:        schema.TypeString,
				Description: "The body of the request",
				Optional:    true,
			},
			// Send Email Action
			"to": {
				Type:        schema.TypeList,
				Description: "The email addresses to which the email is to be sent",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"subject": {
				Type:        schema.TypeString,
				Description: "The subject of the email",
				Optional:    true,
			},
			// body is needed for email as well
			// Trigger Manual Webhook Action
			"webhook_id": {
				Type:        schema.TypeString,
				Description: "The ID of the webhook to be triggered. (Only for Trigger Manual Webhook action)",
				Optional:    true,
			},
			// Status Page Issue Action
			"status_page_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the status page to which the issue is to be added. (Only for Add Status Page Issue action)",
				Optional:    true,
			},
			"issue_title": {
				Type:        schema.TypeString,
				Description: "The title of the issue to be added. (Only for Add Status Page Issue action)",
				Optional:    true,
			},
			"page_status_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the status to be set for the issue. (Only for Add Status Page Issue action)",
				Optional:    true,
			},
			"component_and_impact": {
				Type:        schema.TypeList,
				Description: "The components and their impact to be set for the issue. (Only for Add Status Page Issue action)",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"component_id": {
							Type:        schema.TypeInt,
							Description: "The ID of the component",
							Required:    true,
						},
						"impact_status_id": {
							Type:        schema.TypeInt,
							Description: "The ID of the impact status",
							Required:    true,
						},
					},
				},
			},
			"status_and_message": {
				Type:        schema.TypeList,
				Description: "The status and message to be set for the issue. (Only for Add Status Page Issue action)",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status_id": {
							Type:        schema.TypeInt,
							Description: "The ID of the status",
							Required:    true,
						},
						"messages": {
							Type:        schema.TypeList,
							Description: "The messages to be set for the issue",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			// Jira Create Ticket Action
			"account": {
				Type:        schema.TypeString,
				Description: "The account to be used for creating the ticket. (Only for Jira Create Ticket action)",
				Optional:    true,
			},
			"project": {
				Type:        schema.TypeString,
				Description: "The project to be used for creating the ticket. (Only for Jira Create Ticket action)",
				Optional:    true,
			},
			"issue_type": {
				Type:        schema.TypeString,
				Description: "The issue type to be used for creating the ticket. (Only for Jira Create Ticket action)",
				Optional:    true,
			},
			"title": {
				Type:        schema.TypeString,
				Description: "The title of the ticket. (Only for Jira Create Ticket action)",
				Optional:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the ticket. (Only for Jira Create Ticket action)",
				Optional:    true,
			},
			// Slack:  Channel Creation Action
			"auto_name": {
				Type:        schema.TypeBool,
				Description: "Whether to automatically name the action",
				Optional:    true,
			},
			"channel_name": {
				Type:        schema.TypeString,
				Description: "The name of the channel to be archived. (Only for Slack Archive Channel action)",
				Optional:    true,
			},
			// Slack: Send message to channel
			"channel_id": {
				Type:        schema.TypeString,
				Description: "The ID of the channel to which the message is to be sent. (Only for Slack Message Channel action)",
				Optional:    true,
			},
			"message": {
				Type:        schema.TypeString,
				Description: "The message to be sent. (Only for Slack Message Channel action)",
				Optional:    true,
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
	headers := make([]api.Headers, 0)
	componentAndImpact := make([]api.ComponentAndImpact, 0)
	statusAndMessage := make([]api.StatusAndMessage, 0)

	if err := Decode(d.Get("channels"), &channels); err != nil {
		return diag.FromErr(err)
	}
	if err := Decode(d.Get("headers"), &headers); err != nil {
		return diag.FromErr(err)
	}
	if err := Decode(d.Get("component_and_impact"), &componentAndImpact); err != nil {
		return diag.FromErr(err)
	}
	if err := Decode(d.Get("status_and_message"), &statusAndMessage); err != nil {
		return diag.FromErr(err)
	}

	workflowAction := &api.WorkflowAction{
		Name: d.Get("name").(string),
		Data: api.WorkflowActionData{
			Note:               d.Get("note").(string),
			SLO:                d.Get("slo").(int),
			SLIs:               tf.ListToSlice[string](d.Get("slis")),
			Priority:           d.Get("priority").(string),
			Runbooks:           runbooks,
			Channels:           channels,
			Method:             d.Get("method").(string),
			URL:                d.Get("url").(string),
			Body:               d.Get("body").(string),
			Headers:            headers,
			To:                 tf.ListToSlice[string](d.Get("to")),
			Subject:            d.Get("subject").(string),
			WebhookID:          d.Get("webhook_id").(string),
			StatusPageID:       d.Get("status_page_id").(int),
			IssueTitle:         d.Get("issue_title").(string),
			PageStatusID:       d.Get("page_status_id").(int),
			ComponentAndImpact: componentAndImpact,
			StatusAndMessage:   statusAndMessage,
			Account:            d.Get("account").(string),
			Project:            d.Get("project").(string),
			IssueType:          d.Get("issue_type").(string),
			Title:              d.Get("title").(string),
			Description:        d.Get("description").(string),
			AutoName:           d.Get("auto_name").(bool),
			ChannelName:        d.Get("channel_name").(string),
			ChannelID:          d.Get("channel_id").(string),
			Message:            d.Get("message").(string),
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
	headers := make([]api.Headers, 0)
	componentAndImpact := make([]api.ComponentAndImpact, 0)
	statusAndMessage := make([]api.StatusAndMessage, 0)

	if err := Decode(d.Get("channels"), &channels); err != nil {
		return diag.FromErr(err)
	}
	if err := Decode(d.Get("headers"), &headers); err != nil {
		return diag.FromErr(err)
	}
	if err := Decode(d.Get("component_and_impact"), &componentAndImpact); err != nil {
		return diag.FromErr(err)
	}
	if err := Decode(d.Get("status_and_message"), &statusAndMessage); err != nil {
		return diag.FromErr(err)
	}

	workflowAction := &api.WorkflowAction{
		Name: d.Get("name").(string),
		Data: api.WorkflowActionData{
			Note:               d.Get("note").(string),
			SLO:                d.Get("slo").(int),
			SLIs:               tf.ListToSlice[string](d.Get("slis")),
			Priority:           d.Get("priority").(string),
			Runbooks:           runbooks,
			Channels:           channels,
			Method:             d.Get("method").(string),
			URL:                d.Get("url").(string),
			Body:               d.Get("body").(string),
			Headers:            headers,
			To:                 tf.ListToSlice[string](d.Get("to")),
			Subject:            d.Get("subject").(string),
			WebhookID:          d.Get("webhook_id").(string),
			StatusPageID:       d.Get("status_page_id").(int),
			IssueTitle:         d.Get("issue_title").(string),
			PageStatusID:       d.Get("page_status_id").(int),
			ComponentAndImpact: componentAndImpact,
			StatusAndMessage:   statusAndMessage,
			Account:            d.Get("account").(string),
			Project:            d.Get("project").(string),
			IssueType:          d.Get("issue_type").(string),
			Title:              d.Get("title").(string),
			Description:        d.Get("description").(string),
			AutoName:           d.Get("auto_name").(bool),
			ChannelName:        d.Get("channel_name").(string),
			ChannelID:          d.Get("channel_id").(string),
			Message:            d.Get("message").(string),
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
