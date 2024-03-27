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

func resourceWorkflows() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowsCreate,
		ReadContext:   resourceWorkflowsRead,
		UpdateContext: resourceWorkflowsUpdate,
		DeleteContext: resourceWorkflowsDelete,
		Schema: map[string]*schema.Schema{
			"owner_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
			},
			"title": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 150),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 150),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"trigger": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"incident_created", "incident_triggered", "incident_acknowledged", "incident_resolved"}, false),
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fields": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"color": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"entity_owner": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tf.ValidateObjectID,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"user", "squad", "team"}, false),
						},
					},
				},
			},
		},
	}
}

func resourceWorkflowsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	tflog.Info(ctx, "Damns a new workflow", tf.M{
		"title": d.Get("title").(string),
	})
	client := meta.(*api.Client)

	tflog.Info(ctx, "Creating a new workflow", tf.M{
		"title": d.Get("title").(string),
	})

	// Create a new workflowReq object
	workflowReq := api.Workflow{
		Title:       d.Get("title").(string),
		Description: d.Get("description").(string),
		OwnerID:     d.Get("owner_id").(string),
		Enabled:     d.Get("enabled").(bool),
		Trigger:     d.Get("trigger").(string),
		EntityOwner: api.EntityOwner{
			ID:   d.Get("entity_owner.0.id").(string),
			Type: d.Get("entity_owner.0.type").(string),
		},

		// Filters: []api.Filters{
		// 	{
		// 		Type: d.Get("filters.0.type").(string),
		// 		Fields: api.Field{
		// 			Value: d.Get("filters.0.fields.0.value").(string),
		// 		},
		// 	},
		// },
	}

	mtags := d.Get("tags").([]any)
	tflog.Info(ctx, "Received tags are", tf.M{
		"tags1": mtags,
	})

	if len(mtags) > 0 {
		var tags []api.Tag
		err := Decode(mtags, &tags)
		if err != nil {
			return diag.FromErr(err)
		}
		workflowReq.Tags = tags
	}

	tflog.Info(ctx, "Atleast the basic init is done", tf.M{
		"title": d.Get("title").(string),
	})

	workflow, err := client.CreateWorkflow(ctx, &workflowReq)
	if err != nil {
		return diag.FromErr(err)
	}

	workflowID := strconv.FormatUint(uint64(workflow.ID), 10)
	d.SetId(workflowID)

	return nil // change this
}

func resourceWorkflowsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceWorkflowsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceWorkflowsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
