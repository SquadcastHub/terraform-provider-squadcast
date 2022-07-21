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

func resourceRunbook() *schema.Resource {
	return &schema.Resource{
		Description: "[Squadcast Runbook](https://support.squadcast.com/docs/runbooks) is a compilation of routine procedures and operations that are documented for reference while working on a critical incident. Sometimes, it can also be referred to as a Playbook.",

		CreateContext: resourceRunbookCreate,
		ReadContext:   resourceRunbookRead,
		UpdateContext: resourceRunbookUpdate,
		DeleteContext: resourceRunbookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRunbookImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Runbook id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:  "Name of the Runbook.",
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
			"steps": {
				Description: "Step by Step instructions, you can add as many steps as you want, supports markdown formatting.",
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceRunbookImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*api.Client)

	teamID, name, err := parse2PartImportID(d.Id())
	if err != nil {
		return nil, err
	}

	runbook, err := client.GetRunbookByName(ctx, teamID, name)
	if err != nil {
		return nil, err
	}

	d.Set("team_id", teamID)
	d.SetId(runbook.ID)

	return []*schema.ResourceData{d}, nil
}

func resourceRunbookCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	var steps []*api.RunbookStep
	err := Decode(d.Get("steps"), &steps)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "Creating runbook", tf.M{
		"name": d.Get("name").(string),
	})
	runbook, err := client.CreateRunbook(ctx, &api.CreateUpdateRunbookReq{
		Name:   d.Get("name").(string),
		TeamID: d.Get("team_id").(string),
		Steps:  steps,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(runbook.ID)

	return resourceRunbookRead(ctx, d, meta)
}

func resourceRunbookRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	id := d.Id()

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	tflog.Info(ctx, "Reading runbook", tf.M{
		"id":   d.Id(),
		"name": d.Get("name").(string),
	})
	runbook, err := client.GetRunbookById(ctx, teamID.(string), id)
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(runbook, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRunbookUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	var steps []*api.RunbookStep
	err := Decode(d.Get("steps"), &steps)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateRunbook(ctx, d.Id(), &api.CreateUpdateRunbookReq{
		Name:   d.Get("name").(string),
		TeamID: d.Get("team_id").(string),
		Steps:  steps,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRunbookRead(ctx, d, meta)
}

func resourceRunbookDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteRunbook(ctx, d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
