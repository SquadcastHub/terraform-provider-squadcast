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

func resourceSchedule() *schema.Resource {
	return &schema.Resource{
		Description: "[Squadcast schedules](https://support.squadcast.com/docs/schedules) are used to manage on-call scheduling & determine who will be notified when an incident is triggered.",

		CreateContext: resourceScheduleCreate,
		ReadContext:   resourceScheduleRead,
		UpdateContext: resourceScheduleUpdate,
		DeleteContext: resourceScheduleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceScheduleImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Schedule id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:  "Name of the Schedule.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": {
				Description:  "Detailed description about the Schedule.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"team_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
				ForceNew:     true,
			},
			"color": {
				Description: "Calendar color scheme for this schedule, hex values.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceScheduleImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*api.Client)

	teamID, name, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}

	schedule, err := client.GetScheduleByName(ctx, teamID, name)
	if err != nil {
		return nil, err
	}

	d.Set("team_id", teamID)
	d.SetId(schedule.ID)

	return []*schema.ResourceData{d}, nil
}

func resourceScheduleCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Creating schedule", tf.M{
		"name": d.Get("name").(string),
	})
	schedule, err := client.CreateSchedule(ctx, &api.CreateUpdateScheduleReq{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		TeamID:      d.Get("team_id").(string),
		Color:       d.Get("color").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(schedule.ID)

	return resourceScheduleRead(ctx, d, meta)
}

func resourceScheduleRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	id := d.Id()

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	tflog.Info(ctx, "Reading schedule", tf.M{
		"id":   d.Id(),
		"name": d.Get("name").(string),
	})
	schedule, err := client.GetScheduleById(ctx, teamID.(string), id)
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(schedule, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceScheduleUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.UpdateSchedule(ctx, d.Id(), &api.CreateUpdateScheduleReq{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		TeamID:      d.Get("team_id").(string),
		Color:       d.Get("color").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScheduleRead(ctx, d, meta)
}

func resourceScheduleDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteSchedule(ctx, d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
