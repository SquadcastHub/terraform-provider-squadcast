package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func dataSourceWebform() *schema.Resource {
	return &schema.Resource{
		Description: "[Squadcast Webforms](https://support.squadcast.com/webforms/webforms) allows organizations to expand their customer support by hosting public Webforms, so their customers can quickly create an alert from outside the Squadcast ecosystem. Not only this, but internal stakeholders can also leverage Webforms for easy alert creation. " +
			"Use this data source to get information about a specific webform.",
		ReadContext: dataSourceWebformRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Webform id.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"name": {
				Description: "Name of the Webform.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"team_id": {
				Description:  "Team id.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tf.ValidateObjectID,
			},
			"custom_domain_name": {
				Description: "Custom domain name (URL).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"public_url": {
				Description: "Public URL of the Webform.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"owner": {
				Description: "Form owner.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Form owner type (user, team, squad).",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"id": {
							Description: "Form owner id.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Form owner name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"header": {
				Description: "Webform header.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"title": {
				Description: "Webform title (public).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Description of the Webform.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"footer_text": {
				Description: "Footer text.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"footer_link": {
				Description: "Footer link.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email_on": {
				Description: "Defines when to send email to the reporter (triggered, acknowledged, resolved).",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Description: "Webform Tags.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"services": {
				Description: "Services added to Webform.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Description: "Service ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Service name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"alias": {
							Description: "Service alias.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"input_field": {
				Description: "Input Fields added to Webforms. Added as tags to incident based on selection.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Description: "Input field Label.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"options": {
							Description: "Input field options.",
							Type:        schema.TypeList,
							Computed:    true,
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

func dataSourceWebformRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	name := d.Get("name").(string)

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}

	tflog.Info(ctx, "Reading webform by name", tf.M{
		"name": name,
	})

	webform, err := client.GetWebformByName(ctx, teamID.(string), name)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(webform, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
