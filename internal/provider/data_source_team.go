package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func dataSourceTeam() *schema.Resource {
	return &schema.Resource{
		Description: "Teams is a way for customers to represent their organizational structure in Squadcast. Each Team can be considered as an isolated workspace with their own configurations and permissions." +
			"Use this data source to get information about a specific Team.",
		ReadContext: dataSourceTeamRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Team id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Name of the Team.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Detailed description about the Team.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"default": {
				Description: "Squadcast has one default team and this field let's us know if this is the default team.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Description: "User id.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"role_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Role id.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Role name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"default": {
							Description: "Role is default.",
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
				},
			},
		},
	}
}

func dataSourceTeamRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	name, ok := d.GetOk("name")
	if !ok {
		return diag.Errorf("invalid team name provided")
	}

	team, err := client.GetTeamByName(ctx, name.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(team, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
