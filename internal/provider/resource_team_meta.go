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

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage the Team meta details like Name, descripton etc.",

		CreateContext: resourceTeamCreate,
		ReadContext:   resourceTeamRead,
		UpdateContext: resourceTeamUpdate,
		DeleteContext: resourceTeamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTeamImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Team id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:  "Team name.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"description": {
				Description:  "Team description.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"default": {
				Description: "Team is default?.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"default_role_ids": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTeamImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*api.Client)

	team, err := client.GetTeamByName(ctx, d.Id())
	if err != nil {
		return nil, err
	}
	d.SetId(team.ID)

	return []*schema.ResourceData{d}, nil
}

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Creating team", tf.M{
		"name": d.Get("name").(string),
	})
	team, err := client.CreateTeam(ctx, &api.CreateTeamReq{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(team.ID)

	return resourceTeamRead(ctx, d, meta)
}

func resourceTeamRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Reading team", tf.M{
		"id":   d.Id(),
		"name": d.Get("name").(string),
	})
	team, err := client.GetTeamMetaById(ctx, d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(team, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.UpdateTeamMeta(ctx, d.Id(), &api.UpdateTeamMetaReq{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTeamRead(ctx, d, meta)
}

func resourceTeamDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteTeam(ctx, d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	return nil
}
