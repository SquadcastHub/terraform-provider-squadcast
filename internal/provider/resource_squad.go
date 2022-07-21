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
		Description: "[Squads](https://support.squadcast.com/docs/squads) are smaller groups of members within Teams. Squads could correspond to groups of people that are responsible for specific projects within a Team.",

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
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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

	tflog.Info(ctx, "Creating squad", tf.M{
		"name": d.Get("name").(string),
	})
	squad, err := client.CreateSquad(ctx, &api.CreateSquadReq{
		Name:      d.Get("name").(string),
		MemberIDs: tf.ListToSlice[string](d.Get("member_ids")),
		TeamID:    d.Get("team_id").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(squad.ID)

	return resourceSquadRead(ctx, d, meta)
}

func resourceSquadRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	id := d.Id()

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	tflog.Info(ctx, "Reading squad", tf.M{
		"id":   d.Id(),
		"name": d.Get("name").(string),
	})
	squad, err := client.GetSquadById(ctx, teamID.(string), id)
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

	_, err := client.UpdateSquad(ctx, d.Id(), &api.UpdateSquadReq{
		Name:      d.Get("name").(string),
		MemberIDs: tf.ListToSlice[string](d.Get("member_ids")),
	})
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
