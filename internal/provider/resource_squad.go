package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func resourceSquad() *schema.Resource {
	return &schema.Resource{
		Description: "[Squads](https://support.squadcast.com/docs/squads) are smaller groups of members within Teams. Squads could correspond to groups of people that are responsible for specific projects within a Team. The name of the Squad must be unique within a Team.",

		CreateContext: resourceSquadCreate,
		ReadContext:   resourceSquadRead,
		UpdateContext: resourceSquadUpdate,
		DeleteContext: resourceSquadDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSquadImport,
		},

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
				ForceNew:     true,
			},
			"member_ids": {
				Description: "User ObjectId.",
				Type:        schema.TypeList,
				Deprecated:  "Use `members` instead of `member_ids`.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"members": {
				Description: "list of members belonging to this squad",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Description: "user id.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"role": {
							Description:  "Role of the user. Supported values are 'owner' or 'member' (pass this if your org is using OBAC permission model)",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"owner", "member"}, false),
						},
					},
				},
			},
		},
	}
}

func parse2PartImportID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of import resource id (%s), expected teamID:ID", id)
	}

	return parts[0], parts[1], nil
}

func resourceSquadImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	teamID, id, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("team_id", teamID)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func resourceSquadCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)
	createReq := &api.CreateSquadReq{
		Name:   d.Get("name").(string),
		TeamID: d.Get("team_id").(string),
	}
	memberIDs := tf.ListToSlice[string](d.Get("member_ids"))
	members := d.Get("members").([]interface{})

	if len(members) > 0 && len(memberIDs) > 0 {
		return diag.Errorf("member_ids and members cannot be passed at once")
	}

	if len(memberIDs) > 0 {
		membersArr := make([]api.Member, 0)
		for _, memberID := range memberIDs {
			membersArr = append(membersArr, api.Member{
				UserID: memberID,
				Role:   "member",
			})
		}
		createReq.Members = membersArr
	}

	if len(members) > 0 {
		membersArr := make([]api.Member, 0)
		for _, member := range members {
			mem, ok := member.(map[string]interface{})
			if !ok {
				return diag.Errorf("invalid member")
			}
			membersArr = append(membersArr, api.Member{
				UserID: mem["user_id"].(string),
				Role:   mem["role"].(string),
			})
		}
		createReq.Members = membersArr
	}

	tflog.Info(ctx, "Creating squad", tf.M{
		"name": d.Get("name").(string),
	})
	squad, err := client.CreateSquad(ctx, createReq)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(squad.ID)

	return resourceSquadRead(ctx, d, meta)
}

func resourceSquadRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	id := d.Id()

	tflog.Info(ctx, "Reading squad", tf.M{
		"id":   d.Id(),
		"name": d.Get("name").(string),
	})
	squad, err := client.GetSquadById(ctx, id)
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(squad, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSquadUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	updateReq := &api.UpdateSquadReq{
		Name: d.Get("name").(string),
	}
	memberIDs := tf.ListToSlice[string](d.Get("member_ids"))
	members := d.Get("members").([]interface{})

	// if len(members) > 0 && len(memberIDs) > 0 {
	// 	return diag.Errorf("member_ids and members cannot be passed at once")
	// }

	if len(members) > 0 {
		membersArr := make([]api.Member, 0)
		for _, member := range members {
			mem, ok := member.(map[string]interface{})
			if !ok {
				return diag.Errorf("invalid member")
			}
			membersArr = append(membersArr, api.Member{
				UserID: mem["user_id"].(string),
				Role:   mem["role"].(string),
			})
		}
		updateReq.Members = membersArr
	}

	if len(memberIDs) > 0 {
		membersArr := make([]api.Member, 0)
		for _, memberID := range memberIDs {
			membersArr = append(membersArr, api.Member{
				UserID: memberID,
				Role:   "member",
			})
		}
		updateReq.Members = membersArr
	}

	_, err := client.UpdateSquad(ctx, d.Id(), updateReq)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSquadRead(ctx, d, meta)
}

func resourceSquadDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteSquad(ctx, d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
