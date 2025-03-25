package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func resourceDedupKeyOverlay() *schema.Resource {
	return &schema.Resource{
		Description:   "Define [dedup keys](https://support.squadcast.com/services/alert-deduplication-rules/key-based-deduplication) using customizable templates for configured alert sources.",
		CreateContext: resourceDedupKeyOverlayCreate,
		ReadContext:   resourceDedupKeyOverlayRead,
		UpdateContext: resourceDedupKeyOverlayUpdate,
		DeleteContext: resourceDedupKeyOverlayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDedupKeyOverlayImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of this resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"service_id": {
				Description:  "Service id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
			},
			"alert_source": {
				Description: "Alert source for which the template is defined. Find all alert sources supported on Squadcast [here](https://www.squadcast.com/integrations).",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"alert_source_shortname": {
				Description: "Shortname of the linked alert source.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"dedup_key_overlay_template": {
				Description: "Template for the incident message.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"duration": {
				Description: "Deduplication Time (in minutes). Maximum value is 2880 minutes (48 hours).",
				Type:        schema.TypeInt,
				Required:    true,
			},
		},
	}
}

func resourceDedupKeyOverlayImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	serviceID, alertSourceName, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}

	alertSource, err := api.GetAlertSourceDetailsByName(meta.(*api.Client), ctx, alertSourceName)
	if err != nil {
		return nil, err
	}
	d.Set("service_id", serviceID)
	d.Set("alert_source_shortname", alertSource.ShortName)
	d.Set("alert_source", alertSourceName)
	d.SetId(serviceID)

	return []*schema.ResourceData{d}, nil
}

func resourceDedupKeyOverlayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	serviceID := d.Get("service_id").(string)

	tflog.Info(ctx, "Creating a new dedup key overlay", tf.M{
		"service_id":   serviceID,
		"alert_source": d.Get("alert_source").(string),
	})

	req := api.DedupKeyOverlayReq{
		TemplateType: "go",
		DedupKeyOverlay: api.DedupKey{
			Template:       d.Get("dedup_key_overlay_template").(string),
			DurationInMins: uint(d.Get("duration").(int)),
		},
	}

	alertSourceName := d.Get("alert_source").(string)
	alertSource, err := api.GetAlertSourceDetailsByName(client, ctx, alertSourceName)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.CreateOrUpdateDedupKeyOverlay(ctx, serviceID, alertSource.ShortName, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serviceID)
	d.Set("alert_source_shortname", alertSource.ShortName)

	return resourceDedupKeyOverlayRead(ctx, d, meta)
}

func resourceDedupKeyOverlayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Reading dedup key overlay", tf.M{
		"id":                     d.Id(),
		"alert_source_shortname": d.Get("alert_source_shortname").(string),
	})
	dedupKeyOverlay, err := client.GetDedupKeyOverlay(ctx, d.Id(), d.Get("alert_source_shortname").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(dedupKeyOverlay, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDedupKeyOverlayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)
	serviceID := d.Get("service_id").(string)
	tflog.Info(ctx, "Updating dedup key overlay", tf.M{
		"service_id":   serviceID,
		"alert_source": d.Get("alert_source").(string),
	})

	req := api.DedupKeyOverlayReq{
		TemplateType: "go",
		DedupKeyOverlay: api.DedupKey{
			Template:       d.Get("dedup_key_overlay_template").(string),
			DurationInMins: uint(d.Get("duration").(int)),
		},
	}

	alertSourceName := d.Get("alert_source").(string)
	alertSource, err := api.GetAlertSourceDetailsByName(client, ctx, alertSourceName)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.CreateOrUpdateDedupKeyOverlay(ctx, serviceID, alertSource.ShortName, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("alert_source_shortname", alertSource.ShortName)

	return resourceDedupKeyOverlayRead(ctx, d, meta)
}

func resourceDedupKeyOverlayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Deleting dedup key overlay", tf.M{
		"id":                     d.Id(),
		"alert_source_shortname": d.Get("alert_source_shortname").(string),
	})

	err := client.DeleteDedupKeyOverlay(ctx, d.Id(), d.Get("alert_source_shortname").(string))
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
