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

func dataSourceRunbook() *schema.Resource {
	return &schema.Resource{
		Description: "A Runbook is a compilation of routine procedures and operations that are documented for reference while working on a critical incident. Sometimes, it can also be referred to as a Playbook." +

			"Use this data source to get information about a specific Runbook that you can use for other Squadcast resources.",
		ReadContext: dataSourceRunbookRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Runbook id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:  "Name of the Runbook",
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
			"steps": {
				Description: "Step by Step instructions, you can add as many steps as you want, supports markdown formatting.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRunbookRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	name, ok := d.GetOk("name")
	if !ok {
		return diag.Errorf("invalid runbook name provided")
	}

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	tflog.Info(ctx, "Reading runbook by name", tf.M{
		"name": name.(string),
	})
	runbook, err := client.GetRunbookByName(ctx, teamID.(string), name.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(runbook, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
