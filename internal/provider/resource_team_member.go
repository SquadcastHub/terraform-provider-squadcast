package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func resourceTeamMember() *schema.Resource {
	return &schema.Resource{
		Description: "You can manage the members of a Team here.",

		CreateContext: resourceTeamMemberCreate,
		ReadContext:   resourceTeamMemberRead,
		UpdateContext: resourceTeamMemberUpdate,
		DeleteContext: resourceTeamMemberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTeamMemberImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "id.",
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
			"user_id": {
				Description:  "user id?.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
				ForceNew:     true,
			},
			"role_ids": {
				Description: "role ids.",
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: tf.ValidateObjectID,
				},
			},
		},
	}
}

func resourceTeamMemberImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*api.Client)
	teamID, email, err := parse2PartImportID(d.Id())

	_, err = client.GetTeamById(ctx, teamID)
	if err != nil {
		return nil, err
	}

	user, err := client.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	d.Set("team_id", teamID)
	d.SetId(user.ID)

	return []*schema.ResourceData{d}, nil
}

func resourceTeamMemberCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	teamMember, err := client.CreateTeamMember(ctx, d.Get("team_id").(string), &api.CreateTeamMemberReq{
		UserID:  d.Get("user_id").(string),
		RoleIDs: tf.ListToSlice[string](d.Get("role_ids")),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(teamMember.UserID)

	return resourceTeamMemberRead(ctx, d, meta)
}

func resourceTeamMemberRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	teamMember, err := client.GetTeamMemberByID(ctx, d.Get("team_id").(string), d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(teamMember, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceTeamMemberUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.UpdateTeamMember(ctx, d.Get("team_id").(string), d.Id(), &api.UpdateTeamMemberReq{
		RoleIDs: tf.ListToSlice[string](d.Get("role_ids")),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTeamMemberRead(ctx, d, meta)
}

func resourceTeamMemberDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteTeamMember(ctx, d.Get("team_id").(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
