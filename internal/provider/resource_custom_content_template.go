package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func resourceCustomContentTemplate() *schema.Resource {
	return &schema.Resource{
		Description:   "[Custom Content Templates](https://support.squadcast.com/services/custom-content-templates) empower users to define personalized incident message and description templates by utilizing the payload of a configured alert source for this service.",
		CreateContext: resourceCustomContentTemplateCreate,
		ReadContext:   resourceCustomContentTemplateRead,
		UpdateContext: resourceCustomContentTemplateUpdate,
		DeleteContext: resourceCustomContentTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCustomContentTemplateImport,
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
			"message_template": {
				Description: "Template for the incident message.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description_template": {
				Description: "Template for the incident description.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceCustomContentTemplateImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
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

func resourceCustomContentTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	serviceID := d.Get("service_id").(string)

	tflog.Info(ctx, "Creating a new custom content template", tf.M{
		"service_id":   serviceID,
		"alert_source": d.Get("alert_source").(string),
	})

	req := api.CustomContentTemplateOverlayReq{
		TemplateType: "go",
		MessageOverlay: api.OverlayTemplate{
			Template: d.Get("message_template").(string),
		},
		DescriptionOverlay: api.OverlayTemplate{
			Template: d.Get("description_template").(string),
		},
	}

	alertSourceName := d.Get("alert_source").(string)
	alertSource, err := api.GetAlertSourceDetailsByName(client, ctx, alertSourceName)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.CreateOrUpdateCustomContentTemplateOverlay(ctx, serviceID, alertSource.ShortName, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serviceID)
	d.Set("alert_source_shortname", alertSource.ShortName)

	return resourceCustomContentTemplateRead(ctx, d, meta)
}

func resourceCustomContentTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Reading custom content template", tf.M{
		"id":                     d.Id(),
		"alert_source_shortname": d.Get("alert_source_shortname").(string),
	})
	customContentTemplate, err := client.GetCustomContentTemplateOverlay(ctx, d.Id(), d.Get("alert_source_shortname").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(customContentTemplate, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCustomContentTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)
	serviceID := d.Get("service_id").(string)
	tflog.Info(ctx, "Updating custom content template", tf.M{
		"service_id":   serviceID,
		"alert_source": d.Get("alert_source").(string),
	})

	req := api.CustomContentTemplateOverlayReq{
		TemplateType: "go",
		MessageOverlay: api.OverlayTemplate{
			Template: d.Get("message_template").(string),
		},
		DescriptionOverlay: api.OverlayTemplate{
			Template: d.Get("description_template").(string),
		},
	}

	alertSourceName := d.Get("alert_source").(string)
	alertSource, err := api.GetAlertSourceDetailsByName(client, ctx, alertSourceName)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.CreateOrUpdateCustomContentTemplateOverlay(ctx, serviceID, alertSource.ShortName, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("alert_source_shortname", alertSource.ShortName)

	return resourceCustomContentTemplateRead(ctx, d, meta)
}

func resourceCustomContentTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Deleting custom content template", tf.M{
		"id":                     d.Id(),
		"alert_source_shortname": d.Get("alert_source_shortname").(string),
	})

	err := client.DeleteCustomContentTemplateOverlay(ctx, d.Id(), d.Get("alert_source_shortname").(string))
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
