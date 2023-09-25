package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func resourceGER() *schema.Resource {
	return &schema.Resource{
		Description: "Global Event Ruleset (GER) is a centralized set of rules that defines service routes for incoming events.",

		CreateContext: resourceGERCreate,
		ReadContext:   resourceGERRead,
		UpdateContext: resourceGERUpdate,
		DeleteContext: resourceGERDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGERImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "GER id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"team_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
			},
			"name": {
				Description:  "GER name.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": {
				Description: "GER description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"routing_key": {
				Description: "Routing Key is an identifier used to determine the ruleset that an incoming event belongs to. It is a common key that associates multiple alert sources with their configured rules, ensuring events are routed to the appropriate services when the defined criteria are met.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"entity_owner": {
				Description: "GER owner.",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description:  "GER owner id.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tf.ValidateObjectID,
						},
						"type": {
							Description:  "GER owner type. (user or squad or team)",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"user", "squad", "team"}, false),
						},
					},
				},
			},
		},
	}
}

func resourceGERImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	d.SetId(d.Id())

	return []*schema.ResourceData{d}, nil
}

func resourceGERCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	req := &api.GER{
		TeamID:      d.Get("team_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	entityOwner := d.Get("entity_owner").([]interface{})
	if len(entityOwner) > 0 {
		entityOwnerMap, ok := entityOwner[0].(map[string]interface{})
		if !ok {
			return diag.Errorf("invalid entity_owner")
		}
		req.EntityOwner = &api.EntityOwner{
			ID:   entityOwnerMap["id"].(string),
			Type: entityOwnerMap["type"].(string),
		}
	}

	tflog.Info(ctx, "Creating GER", tf.M{
		"name": req.Name,
	})
	ger, err := client.CreateGER(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	gerID := strconv.FormatUint(uint64(ger.ID), 10)
	d.SetId(gerID)

	return resourceGERRead(ctx, d, meta)
}

func resourceGERRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	id := d.Id()

	tflog.Info(ctx, "Reading GER", tf.M{
		"id": id,
	})
	ger, err := client.GetGERById(ctx, id)
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(ger, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGERUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	req := &api.GER{
		TeamID:      d.Get("team_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	entityOwner := d.Get("entity_owner").([]interface{})
	if len(entityOwner) > 0 {
		entityOwnerMap, ok := entityOwner[0].(map[string]interface{})
		if !ok {
			return diag.Errorf("invalid entity_owner")
		}
		req.EntityOwner = &api.EntityOwner{
			ID:   entityOwnerMap["id"].(string),
			Type: entityOwnerMap["type"].(string),
		}
	}

	tflog.Info(ctx, "Updating GER", tf.M{
		"id": d.Id(),
	})
	_, err := client.UpdateGER(ctx, d.Id(), req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGERRead(ctx, d, meta)
}

func resourceGERDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteGER(ctx, d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
