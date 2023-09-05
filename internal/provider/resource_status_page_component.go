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

func resourceStatusPageComponent() *schema.Resource {
	return &schema.Resource{
		Description: "Status page component defines a component that represents a specific element within a status page. This resource enables you to configure various attributes of the component, and optionally associate it with a group on the status page.",

		CreateContext: resourceStatusPageComponentCreate,
		ReadContext:   resourceStatusPageComponentRead,
		UpdateContext: resourceStatusPageComponentUpdate,
		DeleteContext: resourceStatusPageComponentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceStatusPageComponentImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Component id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status_page_id": {
				Description: "Id of the status page to which this component belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description:  "Name of the status page component.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": {
				Description: "Description of the status page component.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"group_id": {
				Description: "Id of the group to which this component belongs to.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceStatusPageComponentImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	pageID, componentID, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}
	d.Set("status_page_id", pageID)
	d.SetId(componentID)

	return []*schema.ResourceData{d}, nil
}

func resourceStatusPageComponentCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	createStatusPageComponentReq := &api.StatusPageComponent{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if d.Get("group_id").(string) != "" {
		groupId, err := strconv.ParseInt(d.Get("group_id").(string), 10, 64)
		if err != nil {
			return diag.FromErr(err)
		}
		groupID := uint(groupId)
		createStatusPageComponentReq.GroupID = &groupID
	}

	spc, err := client.CreateStatusPageComponent(ctx, d.Get("status_page_id").(string), createStatusPageComponentReq)
	if err != nil {
		return diag.FromErr(err)
	}

	id := strconv.FormatUint(uint64(spc.ID), 10)
	d.SetId(id)

	return resourceStatusPageComponentRead(ctx, d, meta)
}

func resourceStatusPageComponentRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Reading statuspage component", tf.M{
		"id": d.Id(),
	})
	pageId := d.Get("status_page_id").(string)

	spc, err := client.GetStatusPageComponentById(ctx, pageId, d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(spc, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceStatusPageComponentUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	updateStatusPageReq := &api.StatusPageComponent{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if d.Get("group_id").(string) == "" {
		flag := false
		updateStatusPageReq.BelongsToGroup = &flag
	} else {
		groupId, err := strconv.ParseInt(d.Get("group_id").(string), 10, 64)
		if err != nil {
			return diag.FromErr(err)
		}
		flag := true
		groupID := uint(groupId)
		updateStatusPageReq.BelongsToGroup = &flag
		updateStatusPageReq.GroupID = &groupID
	}

	_, err := client.UpdateStatusPageComponent(ctx, d.Get("status_page_id").(string), d.Id(), updateStatusPageReq)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStatusPageComponentRead(ctx, d, meta)
}

func resourceStatusPageComponentDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	pageId := d.Get("status_page_id").(string)
	_, err := client.DeleteStatusPageComponent(ctx, pageId, d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
