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

func dataSourceTeamRole() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information about various Team Roles.",

		ReadContext: dataSourceTeamRoleRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Role id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"team_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
				ForceNew:     true,
			},
			"name": {
				Description:  "TeamRole name.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"default": {
				Description: "Role is default?.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"abilities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceTeamRoleRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	teamRoleName := d.Get("name").(string)
	team_id := d.Get("team_id").(string)

	tflog.Info(ctx, "Reading team_role", tf.M{
		"name": teamRoleName,
		"id":   team_id,
	})
	teamRole, err := client.GetTeamRoleByName(ctx, team_id, teamRoleName)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(teamRole, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
