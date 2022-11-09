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

func resourceTeamRole() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage the Team roles and their permissions",

		CreateContext: resourceTeamRoleCreate,
		ReadContext:   resourceTeamRoleRead,
		UpdateContext: resourceTeamRoleUpdate,
		DeleteContext: resourceTeamRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTeamRoleImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "TeamRole id.",
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
				Description:  "Team role name.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"default": {
				Description: "Team role default.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"abilities": {
				Description: "abilities. \n Current available abilities are : \n create-escalation-policies, create-postmortems, create-runbooks, create-schedules, create-services, create-slos, create-squads, create-status-pages, delete-escalation-policies, delete-postmortems, delete-runbooks, delete-schedules, delete-services, delete-slos, delete-squads, delete-status-pages, read-escalation-policies, read-postmortems, read-runbooks, read-schedules, read-services, read-slos, read-squads, read-status-pages, read-team-analytics, update-escalation-policies, update-postmortems, update-runbooks, update-schedules, update-services, update-slos, update-squads, update-status-pages",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTeamRoleImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	teamID, teamRoleName, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}

	client := meta.(*api.Client)

	teamRole, err := client.GetTeamRoleByName(ctx, teamID, teamRoleName)
	if err != nil {
		return nil, err
	}

	d.Set("team_id", teamID)
	d.SetId(teamRole.ID)

	return []*schema.ResourceData{d}, nil
}

func resourceTeamRoleCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Creating team_role", tf.M{
		"name": d.Get("name").(string),
	})
	teamRole, err := client.CreateTeamRole(ctx, d.Get("team_id").(string), &api.CreateTeamRoleReq{
		Name:      d.Get("name").(string),
		Abilities: tf.ListToSlice[string](d.Get("abilities")),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(teamRole.ID)

	return resourceTeamRoleRead(ctx, d, meta)
}

func resourceTeamRoleRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	id := d.Id()

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	tflog.Info(ctx, "Reading team_role", tf.M{
		"id":   d.Id(),
		"name": d.Get("name").(string),
	})
	teamRole, err := client.GetTeamRoleByID(ctx, teamID.(string), id)
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(teamRole, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceTeamRoleUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.UpdateTeamRole(ctx, d.Get("team_id").(string), d.Id(), &api.UpdateTeamRoleReq{
		Name:      d.Get("name").(string),
		Abilities: tf.ListToSlice[string](d.Get("abilities")),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTeamRoleRead(ctx, d, meta)
}

func resourceTeamRoleDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteTeamRole(ctx, d.Get("team_id").(string), d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
