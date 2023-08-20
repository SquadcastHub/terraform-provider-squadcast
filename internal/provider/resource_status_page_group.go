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

func resourceStatusPageGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Status page resource.",

		CreateContext: resourceStatusPageGroupCreate,
		ReadContext:   resourceStatusPageGroupRead,
		UpdateContext: resourceStatusPageGroupUpdate,
		DeleteContext: resourceStatusPageGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceStatusPageGroupImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Group id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status_page_id": {
				Description: "ID of the status page to which this group belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description:  "Status page group name.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": {
				Description: "Description of the status page group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"allow_subscription": {
				Description: "Allow subscription to the status page group.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"component_ids": {
				Description: "List of component ids that belong to this group.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceStatusPageGroupImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*api.Client)
	pageId, groupId, err := parse2PartImportID(d.Id())

	spg, err := client.GetStatusPageGroupById(ctx, pageId, groupId)
	if err != nil {
		return nil, err
	}

	id := strconv.FormatUint(uint64(spg.ID), 10)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func resourceStatusPageGroupCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	createStatusPageGroupReq := &api.StatusPageGroup{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		AllowSubscription: d.Get("allow_subscription").(bool),
	}

	componentIds := d.Get("component_ids").([]interface{})
	for _, componentId := range componentIds {
		id, err := strconv.ParseUint(componentId.(string), 10, 32)
		if err != nil {
			return diag.FromErr(err)
		}
		createStatusPageGroupReq.ComponentIDs = append(createStatusPageGroupReq.ComponentIDs, uint(id))
	}

	spg, err := client.CreateStatusPageGroup(ctx, d.Get("status_page_id").(string), createStatusPageGroupReq)
	if err != nil {
		return diag.FromErr(err)
	}

	id := strconv.FormatUint(uint64(spg.ID), 10)
	d.SetId(id)

	return resourceStatusPageGroupRead(ctx, d, meta)
}

func resourceStatusPageGroupRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Reading statuspage group", tf.M{
		"id": d.Id(),
	})
	pageId := d.Get("status_page_id").(string)

	spg, err := client.GetStatusPageGroupById(ctx, pageId, d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(spg, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceStatusPageGroupUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	updateStatusPageGroupReq := &api.StatusPageGroup{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		AllowSubscription: d.Get("allow_subscription").(bool),
	}

	componentIds := d.Get("component_ids").([]interface{})
	for _, componentId := range componentIds {
		id, err := strconv.ParseUint(componentId.(string), 10, 32)
		if err != nil {
			return diag.FromErr(err)
		}
		updateStatusPageGroupReq.ComponentIDs = append(updateStatusPageGroupReq.ComponentIDs, uint(id))
	}

	_, err := client.UpdateStatusPageGroup(ctx, d.Get("status_page_id").(string), d.Id(), updateStatusPageGroupReq)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStatusPageGroupRead(ctx, d, meta)
}

func resourceStatusPageGroupDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	pageId := d.Get("status_page_id").(string)
	_, err := client.DeleteStatusPageGroup(ctx, pageId, d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
