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

func resourceWorkflow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowsCreate,
		ReadContext:   resourceWorkflowsRead,
		UpdateContext: resourceWorkflowsUpdate,
		DeleteContext: resourceWorkflowsDelete,
		Schema: map[string]*schema.Schema{
			"owner_id": {
				Type:         schema.TypeString,
				Description:  "The ID of the user who owns the workflow",
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
			},
			"title": {
				Type:         schema.TypeString,
				Description:  "The title of the workflow",
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 150),
			},
			"description": {
				Type:         schema.TypeString,
				Description:  "The description of the workflow",
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 150),
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the workflow is enabled or not",
				Optional:    true,
				Default:     true,
			},
			"trigger": {
				Type:         schema.TypeString,
				Description:  "The trigger for the workflow",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"incident_created", "incident_triggered", "incident_acknowledged", "incident_resolved"}, false),
			},
			"filters": {
				Type:        schema.TypeList,
				Description: "The filters to be applied on the workflow",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition": {
							Type:        schema.TypeString,
							Description: "Condition to be applied on the filters (and / or)",
							Required:    true,
						},
						"filters": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"condition": {
										Type:        schema.TypeString,
										Description: "Condition to be applied on the filters (and / or)",
										Optional:    true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"filters": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"key": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"value": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeList,
				Description: "The tags to be applied on the workflow",
				Optional:    true,
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
				Type:        schema.TypeList,
				Description: "The entity owner of the workflow",
				Required:    true,
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

	client := meta.(*api.Client)

	tflog.Info(ctx, "Creating a new workflow", tf.M{
		"title": d.Get("title").(string),
	})

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
	}

	hfilters := d.Get("filters").([]any)
	if len(hfilters) > 0 {
		var filters []*api.HighLevelFilter
		err := Decode(hfilters, &filters)
		if err != nil {
			return diag.FromErr(err)
		}

		workflowReq.Filters = filters[0]
	}

	mtags := d.Get("tags").([]any)
	tflog.Info(ctx, "Length of tags from create are", tf.M{
		"tagsLen": len(mtags),
	})

	if len(mtags) > 0 {
		var tags []*api.WorkflowTag
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

	return resourceWorkflowsRead(ctx, d, meta)
}

func resourceWorkflowsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Reading workflow", tf.M{
		"id": d.Id(),
	})
	workflow, err := client.GetWorkflowById(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "debug: Called getByWorkflowID", tf.M{
		"id": d.Id(),
	})

	if err = tf.EncodeAndSet(workflow, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWorkflowsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Updating workflow", tf.M{
		"id": d.Id(),
	})

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
	}

	hfilters := d.Get("filters").([]any)
	if len(hfilters) > 0 {
		var filters []*api.HighLevelFilter
		err := Decode(hfilters, &filters)
		if err != nil {
			return diag.FromErr(err)
		}

		workflowReq.Filters = filters[0]
	}

	mtags := d.Get("tags").([]any)
	tflog.Info(ctx, "Received tags from update are", tf.M{
		"tags1": mtags,
	})

	if len(mtags) > 0 {
		var tags []*api.WorkflowTag
		err := Decode(mtags, &tags)
		if err != nil {
			return diag.FromErr(err)
		}

		workflowReq.Tags = tags
	}

	_, err := client.UpdateWorkflow(ctx, d.Id(), &workflowReq)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceWorkflowsRead(ctx, d, meta)
}

func resourceWorkflowsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Deleting workflow", tf.M{
		"id": d.Id(),
	})

	_, err := client.DeleteWorkflow(ctx, d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
