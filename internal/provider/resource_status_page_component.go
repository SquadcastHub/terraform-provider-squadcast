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
		Description: "Status page resource.",

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
				Description: "ID of the status page to which this component belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description:  "Status page component name.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": {
				Description: "Description of the status page component.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"allow_subscription": {
				Description: "Allow subscription to the status page component.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceStatusPageComponentImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*api.Client)
	pageId, componentId, err := parse2PartImportID(d.Id())

	spc, err := client.GetStatusPageComponentById(ctx, pageId, componentId)
	if err != nil {
		return nil, err
	}

	id := strconv.FormatUint(uint64(spc.ID), 10)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func resourceStatusPageComponentCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	createStatusPageComponentReq := &api.StatusPageComponent{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		AllowSubscription: d.Get("allow_subscription").(bool),
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
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		AllowSubscription: d.Get("allow_subscription").(bool),
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
