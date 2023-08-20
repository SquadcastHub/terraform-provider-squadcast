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

func resourceStatusPage() *schema.Resource {
	return &schema.Resource{
		Description: "Status page resource.",

		CreateContext: resourceStatusPageCreate,
		ReadContext:   resourceStatusPageRead,
		UpdateContext: resourceStatusPageUpdate,
		DeleteContext: resourceStatusPageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceStatusPageImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Status page id.",
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
				Description:  "Status page name.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": {
				Description: "Description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"is_public": {
				Description: "Description.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"domain_name": {
				Description: "Domain name of the status page. This will be appended to https://statuspage.squadcast.com/<ORG_ID>/ to form the URL of the status page.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"custom_domain_name": {
				Description: "Custom domain name of the status page.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"timezone": {
				Description: "Timezone",
				Type:        schema.TypeString,
				Required:    true,
			},
			"contact_email": {
				Description: "Contact email.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"theme_color": {
				Description: "Theme color.",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"primary": {
							Description: "Primary color.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"secondary": {
							Description: "Secondary color.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"owner": {
				Description: "Status page owner.",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description:  "Status page owner type (user, team, squad).",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"user", "squad", "team"}, false),
						},
						"id": {
							Description:  "Status page owner id.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tf.ValidateObjectID,
						},
					},
				},
			},
			"allow_webhook_subscription": {
				Description: "Allow webhook subscription to the status page.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"allow_maintenance_subscription": {
				Description: "Allow maintenance subscription to the status page.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"allow_components_subscription": {
				Description: "Allow components subscription to the status page.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceStatusPageImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*api.Client)

	sp, err := client.GetStatusPageById(ctx, d.Id())
	if err != nil {
		return nil, err
	}

	id := strconv.FormatUint(uint64(sp.ID), 10)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func resourceStatusPageCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	createStatusPageReq := &api.StatusPage{
		TeamID:                       d.Get("team_id").(string),
		Name:                         d.Get("name").(string),
		Description:                  d.Get("description").(string),
		IsPublic:                     d.Get("is_public").(bool),
		DomainName:                   d.Get("domain_name").(string),
		Timezone:                     d.Get("timezone").(string),
		ContactEmail:                 d.Get("contact_email").(string),
		AllowWebhookSubscription:     d.Get("allow_webhook_subscription").(bool),
		AllowMaintenanceSubscription: d.Get("allow_maintenance_subscription").(bool),
		AllowComponentsSubscription:  d.Get("allow_components_subscription").(bool),
	}

	if d.Get("custom_domain_name").(string) != "" {
		createStatusPageReq.CustomDomainName = d.Get("custom_domain_name").(string)
	}
	ownerData, err := tf.ExtractData(d, "owner")
	if err != nil {
		return diag.FromErr(err)
	}
	createStatusPageReq.OwnerID = ownerData["id"].(string)
	createStatusPageReq.OwnerType = ownerData["type"].(string)

	themeColor, err := tf.ExtractData(d, "theme_color")
	if err != nil {
		return diag.FromErr(err)
	}
	createStatusPageReq.ThemeColor.Primary = themeColor["primary"].(string)
	createStatusPageReq.ThemeColor.Secondary = themeColor["secondary"].(string)

	sp, err := client.CreateStatusPage(ctx, createStatusPageReq)
	if err != nil {
		return diag.FromErr(err)
	}

	id := strconv.FormatUint(uint64(sp.ID), 10)
	d.SetId(id)

	return resourceStatusPageRead(ctx, d, meta)
}

func resourceStatusPageRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	id := d.Id()

	tflog.Info(ctx, "Reading statuspage", tf.M{
		"id": id,
	})
	sp, err := client.GetStatusPageById(ctx, id)
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(sp, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceStatusPageUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	if d.HasChange("team_id") {
		return diag.Errorf("team_id can only be set during creation.")
	}

	updateStatusPageReq := &api.StatusPage{
		TeamID:                       d.Get("team_id").(string),
		Name:                         d.Get("name").(string),
		Description:                  d.Get("description").(string),
		IsPublic:                     d.Get("is_public").(bool),
		DomainName:                   d.Get("domain_name").(string),
		Timezone:                     d.Get("timezone").(string),
		ContactEmail:                 d.Get("contact_email").(string),
		AllowWebhookSubscription:     d.Get("allow_webhook_subscription").(bool),
		AllowMaintenanceSubscription: d.Get("allow_maintenance_subscription").(bool),
		AllowComponentsSubscription:  d.Get("allow_components_subscription").(bool),
	}

	if d.Get("custom_domain_name").(string) != "" {
		updateStatusPageReq.CustomDomainName = d.Get("custom_domain_name").(string)
	}

	ownerData, err := tf.ExtractData(d, "owner")
	if err != nil {
		return diag.FromErr(err)
	}
	updateStatusPageReq.OwnerID = ownerData["id"].(string)
	updateStatusPageReq.OwnerType = ownerData["type"].(string)

	themeColor, err := tf.ExtractData(d, "theme_color")
	if err != nil {
		return diag.FromErr(err)
	}
	updateStatusPageReq.ThemeColor.Primary = themeColor["primary"].(string)
	updateStatusPageReq.ThemeColor.Secondary = themeColor["secondary"].(string)

	_, err = client.UpdateStatusPage(ctx, d.Id(), updateStatusPageReq)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStatusPageRead(ctx, d, meta)
}

func resourceStatusPageDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteStatusPage(ctx, d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
