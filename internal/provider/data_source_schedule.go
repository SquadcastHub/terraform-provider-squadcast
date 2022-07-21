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

func dataSourceSchedule() *schema.Resource {
	return &schema.Resource{
		Description: "[Squadcast schedules](https://support.squadcast.com/docs/schedules) are used to manage on-call scheduling & determine who will be notified when an incident is triggered. " +
			"Use this data source to get information about a specific schedule that you can use for other Squadcast resources.",
		ReadContext: dataSourceScheduleRead,
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
				Description: "Detailed description about the schedule.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"team_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
			},
			"color": {
				Description: "Calendar color scheme for this schedule, hex values.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceScheduleRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	name, ok := d.GetOk("name")
	if !ok {
		return diag.Errorf("invalid schedule name provided")
	}

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	tflog.Info(ctx, "Reading schedule by name", tf.M{
		"name": name.(string),
	})
	schedule, err := client.GetScheduleByName(ctx, teamID.(string), name.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(schedule, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
