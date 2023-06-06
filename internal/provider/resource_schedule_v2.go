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

func resourceScheduleV2() *schema.Resource {
	return &schema.Resource{
		Description:   "Squadcast v2 schedules", //todo: update this
		ReadContext:   resourceScheduleV2Read,
		CreateContext: resourceScheduleV2Create,
		UpdateContext: resourceScheduleV2Create,
		DeleteContext: resourceScheduleV2Delete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Schedule id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:  "Name of the schedule.",
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
			"timezone": {
				Description: "Timezone of the schedule",
				Type:        schema.TypeString,
				Required:    true,
			},
			"entity_owner": {
				Description: "Schedule entity owner.",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description:  "Schedule owner type (user, team, squad).",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"user", "squad", "team"}, false),
						},
						"id": {
							Description:  "Schedule owner id.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tf.ValidateObjectID,
						},
					},
				},
			},
			"team_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
				ForceNew:     true,
			},
		},
	}
}

// func resourceScheduleV2Import(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
// 	client := meta.(*api.Client)

// 	teamID, name, err := parse2PartImportID(d.Id())
// 	if err != nil {
// 		return nil, err
// 	}

// 	schedule, err := client.GetScheduleById(ctx, teamID, name)
// 	if err != nil {
// 		return nil, err
// 	}

// 	d.Set("team_id", teamID)
// 	d.SetId(schedule.ID)

// 	return []*schema.ResourceData{d}, nil
// }

func resourceScheduleV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {

	client := meta.(*api.Client)

	id := d.Id()
	tflog.Info(ctx, "Reading schedule", tf.M{
		"id":   d.Id(),
		"name": d.Get("name").(string),
	})

	schedule, err := client.GetScheduleV2ById(ctx, id)
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

func resourceScheduleV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Creating schedule", tf.M{
		"name": d.Get("name").(string),
	})

	createScheduleReq := api.NewSchedule{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		TimeZone:    d.Get("timezone").(string),
		TeamID:      d.Get("team_id").(string),
	}
	entityOwner := d.Get("entity_owner").([]interface{})
	if len(entityOwner) > 0 {
		entityOwnerMap, ok := entityOwner[0].(map[string]interface{})
		if !ok {
			return diag.Errorf("entity_owner is invalid")
		}
		createScheduleReq.Owner = &api.Owner{
			Type: entityOwnerMap["type"].(string),
			ID:   entityOwnerMap["id"].(string),
		}
	}

	schedule, err := client.CreateScheduleV2(ctx, createScheduleReq)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(schedule.NewSchedule.ID))

	return resourceScheduleV2Read(ctx, d, meta)
}

func resourceScheduleV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteScheduleV2ByID(ctx, d.Id())
	if err != nil {
		tflog.Info(ctx, "No err while deleting schedule")
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			tflog.Info(ctx, "No resource found while deleting schedule")
			return nil
		}
		tflog.Info(ctx, "random err found while deleting schedule")
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "No err while deleting schedule")
	return nil
}
