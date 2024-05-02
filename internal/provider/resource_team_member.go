package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Description:  "user id (ObjectId).",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
				ForceNew:     true,
			},
			"role_ids": {
				Description: "role ids (pass this if your org is using RBAC permission model)",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: tf.ValidateObjectID,
				},
			},
			"role": {
				Description:  "Role of the member. Supported values are 'stakeholder', 'member' or 'owner' (pass this if your org is using OBAC permission model)",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"stakeholder", "owner", "member"}, false),
			},
		},
	}
}

func resourceTeamMemberImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*api.Client)
	teamID, email, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}
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

	createReq := &api.CreateTeamMemberReq{
		UserID: d.Get("user_id").(string),
	}

	roleIDs := tf.ListToSlice[string](d.Get("role_ids"))
	role := d.Get("role").(string)

	if len(roleIDs) > 0 && len(role) > 0 {
		return diag.Errorf("role_ids and role cannot be passed")
	}

	if len(roleIDs) > 0 {
		createReq.RoleIDs = roleIDs
	}

	if len(role) > 0 {
		createReq.Role = role
	}

	teamMember, err := client.CreateTeamMember(ctx, d.Get("team_id").(string), createReq)
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
	updateReq := &api.UpdateTeamMemberReq{}

	roleIDs := tf.ListToSlice[string](d.Get("role_ids"))
	role := d.Get("role").(string)

	if len(roleIDs) > 0 && len(role) > 0 {
		return diag.Errorf("role_ids and role cannot be passed")
	}

	if len(roleIDs) > 0 {
		updateReq.RoleIDs = roleIDs
	}

	if len(role) > 0 {
		updateReq.Role = role
	}

	_, err := client.UpdateTeamMember(ctx, d.Get("team_id").(string), d.Id(), updateReq)
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
