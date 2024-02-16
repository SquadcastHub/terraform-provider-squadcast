package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func dataSourceScheduleV2() *schema.Resource {
	return &schema.Resource{
		Description: "[Squadcast schedules v2](https://support.squadcast.com/schedules/schedules-new) are used to manage on-call scheduling & determine who will be notified when an incident is triggered. " +
			"Use this data source to get information about a specific schedule that you can use for other Squadcast resources.",
		ReadContext: dataSourceScheduleV2Read,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Schedule id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"team_id": {
				Description: "Team id.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Name of the Schedule.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Detailed description about the schedule.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"timezone": {
				Description: "Timezone for the schedule.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"entity_owner": {
				Description: "Schedule owner.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Schedule owner type (user, team, squad).",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"id": {
							Description: "Schedule owner id.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"tags": {
				Description: "Schedule tags.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Description: "Schedule tag key.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"value": {
							Description: "Schedule tag value.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"color": {
							Description: "Schedule tag color.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceScheduleV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	name, ok := d.GetOk("name")
	if !ok {
		return diag.Errorf("invalid schedule name provided")
	}

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team_id provided")
	}

	tflog.Info(ctx, "Reading schedule_v2 by name", tf.M{
		"name":    name.(string),
		"team_id": teamID.(string),
	})

	schedules, err := client.GetScheduleV2ByName(ctx, teamID.(string), name.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(schedules.NewSchedule) == 0 {
		return diag.Errorf("no schedule found with name %s", name.(string))
	}
	schedule := schedules.NewSchedule[0]

	if err = tf.EncodeAndSet(schedule, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
