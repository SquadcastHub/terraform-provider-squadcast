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

func resourceIAGConfig() *schema.Resource {
	return &schema.Resource{
		Description: "[Intelligent Alert Grouping](https://support.squadcast.com/services/intelligent-alert-grouping-iag) automatically group incoming alerts with a similar open incident and save your team from alert noise.",

		CreateContext: resourceIAGConfigCreateOrUpdate,
		ReadContext:   resourceIAGConfigRead,
		UpdateContext: resourceIAGConfigCreateOrUpdate,
		DeleteContext: resourceIAGConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIAGConfigImport,
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
			"grouping_window": {
				Description:  "Grouping window (in minutes).",
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{5, 10, 15, 20, 45, 1 * 60, 2 * 60, 4 * 60, 8 * 60, 12 * 60, 24 * 60}),
			},
		},
	}
}

func resourceIAGConfigImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	d.Set("service_id", d.Id())

	return []*schema.ResourceData{d}, nil
}

func resourceIAGConfigCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)
	serviceID, ok := d.Get("service_id").(string)
	if !ok {
		return diag.Errorf("invalid service_id")
	}

	req := api.IAGConfig{
		IsEnabled:           d.Get("is_enabled").(bool),
		RollingWindowInMins: d.Get("grouping_window").(int),
	}

	tflog.Info(ctx, "Updating IAG Config for service", tf.M{
		"service_id": serviceID,
	})

	_, err := client.UpdateIAGConfig(ctx, serviceID, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serviceID)

	return resourceIAGConfigRead(ctx, d, meta)
}

func resourceIAGConfigRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Reading IAG Config", tf.M{
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

	if err = tf.EncodeAndSet(service.IAGConfig, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceIAGConfigDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.UpdateIAGConfig(ctx, d.Get("service_id").(string), &api.IAGConfig{
		IsEnabled: false,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
