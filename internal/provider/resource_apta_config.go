package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func resourceAPTAConfig() *schema.Resource {
	return &schema.Resource{
		Description: "[Auto Pause Transient Alerts](https://support.squadcast.com/services/auto-pause-transient-alerts-apta) automatically pauses notifications for transient alerts, giving time for them to auto-resolve before notifying responders.",

		CreateContext: resourceAPTAConfigCreateOrUpdate,
		ReadContext:   resourceAPTAConfigRead,
		UpdateContext: resourceAPTAConfigCreateOrUpdate,
		DeleteContext: resourceAPTAConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAPTAConfigImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"service_id": {
				Description:  "Service ID.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
				ForceNew:     true,
			},
			"is_enabled": {
				Description: "Determines whether this setting needs to be enabled or not.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"timeout": {
				Description:  "This is the timeout window (in minutes) for which an alert flagged as transient will remain in the suppressed state. Supported values are 2, 3, 5, 10 and 15.",
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{2, 3, 5, 10, 15}),
			},
		},
	}
}

func resourceAPTAConfigImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	d.Set("service_id", d.Id())

	return []*schema.ResourceData{d}, nil
}

func resourceAPTAConfigCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)
	serviceID, ok := d.Get("service_id").(string)
	if !ok {
		return diag.Errorf("invalid service_id")
	}

	req := api.APTAConfig{
		IsEnabled:     d.Get("is_enabled").(bool),
		TimeoutInMins: d.Get("timeout").(int),
	}

	tflog.Info(ctx, "Updating APTA Config for service", tf.M{
		"service_id": serviceID,
	})

	_, err := client.UpdateAPTAConfig(ctx, serviceID, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serviceID)

	return resourceAPTAConfigRead(ctx, d, meta)
}

func resourceAPTAConfigRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Reading APTA Config", tf.M{
		"service_id": d.Id(),
	})

	service, err := client.GetServiceById(ctx, "", d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(service.APTAConfig, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAPTAConfigDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.UpdateAPTAConfig(ctx, d.Get("service_id").(string), &api.APTAConfig{
		IsEnabled: false,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
