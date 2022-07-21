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

func dataSourceSquad() *schema.Resource {
	return &schema.Resource{
		Description: "[Squads](https://support.squadcast.com/docs/squads) are smaller groups of members within Teams. Squads could correspond to groups of people that are responsible for specific projects within a Team." +
			"Use this data source to get information about a specific Squad.",
		ReadContext: dataSourceSquadRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Squad id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:  "Name of the Squad.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"team_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
			},
			"member_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceSquadRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	name, ok := d.GetOk("name")
	if !ok {
		return diag.Errorf("invalid squad name provided")
	}

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	tflog.Info(ctx, "Reading squad by name", tf.M{
		"name": name.(string),
	})
	squad, err := client.GetSquadByName(ctx, teamID.(string), name.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(squad, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
