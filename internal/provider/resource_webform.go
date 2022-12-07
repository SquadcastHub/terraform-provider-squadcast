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

func resourceWebform() *schema.Resource {
	return &schema.Resource{
		Description: "[Squadcast Webforms](https://support.squadcast.com/webforms/webforms) allows organizations to expand their customer support by hosting public Webforms, so their customers can quickly create an alert from outside the Squadcast ecosystem. Not only this, but internal stakeholders can also leverage Webforms for easy alert creation.",

		CreateContext: resourceWebformCreate,
		ReadContext:   resourceWebformRead,
		UpdateContext: resourceWebformUpdate,
		DeleteContext: resourceWebformDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWebformImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Webform id.",
				Type:        schema.TypeString,
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
				ForceNew:     true,
			},
			"custom_domain_name": {
				Description: "Custom domain name (URL).",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"public_url": {
				Description: "Public URL of the Webform.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"owner": {
				Description: "Form owner.",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems: 	1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Form owner type (user, team, squad).",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "Form owner id.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"name": {
							Description: "Form owner name.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"header": {
				Description: "Webform header.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"title": {
				Description: "Webform title (public).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Description of the Webform.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"footer_text": {
				Description: "Footer text.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"footer_link": {
				Description: "Footer link.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"email_on": {
				Description: "Defines when to send email to the reporter (triggered, acknowledged, resolved).",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"triggered", "acknowledged", "resolved"}, false),
				},
			},
			"tags": {
				Description: "Webform Tags.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"services": {
				Description: "Services added to Webform.",
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Description: "Service ID.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"name": {
							Description: "Service name.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"alias": {
							Description: "Service alias.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"severity": {
				Description: "Severity of the Incident.",
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Severity type.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"description": {
							Description: "Severity description.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceWebformImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*api.Client)
	teamID, webformName, err := parse2PartImportID(d.Id())

	_, err = client.GetTeamById(ctx, teamID)
	if err != nil {
		return nil, err
	}

	webform, err := client.GetWebformByName(ctx, teamID, webformName)
	if err != nil {
		return nil, err
	}
	webform.WebformOwner = &api.WebformOwner{
		Type: webform.FormOwnerType,
		ID:   webform.FormOwnerID,
		Name: webform.FormOwnerName,
	}

	d.Set("team_id", teamID)
	webformId := strconv.FormatUint(uint64(webform.ID), 10)
	d.SetId(webformId)

	return []*schema.ResourceData{d}, nil
}

func resourceWebformCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Creating webform", tf.M{
		"name": d.Get("name").(string),
	})

	webformOwner := d.Get("owner").([]interface{})[0].(map[string]interface{})

	webformCreateReq := api.WebformReq{
		Name:          d.Get("name").(string),
		TeamID:        d.Get("team_id").(string),
		FormOwnerType: webformOwner["type"].(string),
		FormOwnerID:   webformOwner["id"].(string),
		FormOwnerName: webformOwner["name"].(string),
		HostName:      d.Get("custom_domain_name").(string),
		Header:        d.Get("header").(string),
		Description:   d.Get("description").(string),
		Title:         d.Get("title").(string),
		FooterText:    d.Get("footer_text").(string),
		FooterLink:    d.Get("footer_link").(string),
	}

	if d.Get("custom_domain_name").(string) != "" {
		webformCreateReq.IsCname = true
	}

	memailon := d.Get("email_on").([]interface{})
	emailon := make([]string, len(memailon))
	for i, v := range memailon {
		emailon[i] = v.(string)
	}
	webformCreateReq.EmailOn = emailon

	mservices := d.Get("services").([]interface{})

	var services []api.WFService
	err := Decode(mservices, &services)
	if err != nil {
		return diag.FromErr(err)
	}
	webformCreateReq.Services = services

	mseverity := d.Get("severity").([]interface{})
	var severity []api.WFSeverity
	err = Decode(mseverity, &severity)
	if err != nil {
		return diag.FromErr(err)
	}
	webformCreateReq.Severity = severity

	mtags := d.Get("tags").(map[string]interface{})
	tags := make(map[string]string, len(*&mtags))
	for k, v := range *&mtags {
		tags[k] = v.(string)
	}

	webformCreateReq.Tags = tags

	webformRes, err := client.CreateWebform(ctx, d.Get("team_id").(string), &webformCreateReq)
	if err != nil {
		return diag.FromErr(err)
	}
	webform := webformRes.WebFormRes

	webformId := strconv.FormatUint(uint64(webform.ID), 10)
	d.SetId(webformId)

	return resourceWebformRead(ctx, d, meta)
}

func resourceWebformRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	id := d.Id()

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided !")
	}

	tflog.Info(ctx, "Reading webform", tf.M{
		"id":   d.Id(),
		"name": d.Get("name").(string),
	})

	webform, err := client.GetWebformById(ctx, teamID.(string), id)
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(webform, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWebformUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	tflog.Info(ctx, "Creating webform", tf.M{
		"name": d.Get("name").(string),
	})
	webformOwner := d.Get("owner").([]interface{})[0].(map[string]interface{})

	webformUpdateReq := api.WebformReq{
		Name:          d.Get("name").(string),
		TeamID:        d.Get("team_id").(string),
		FormOwnerType: webformOwner["type"].(string),
		FormOwnerID:   webformOwner["id"].(string),
		FormOwnerName: webformOwner["name"].(string),
		HostName:      d.Get("custom_domain_name").(string),
		Header:        d.Get("header").(string),
		Description:   d.Get("description").(string),
		Title:         d.Get("title").(string),
		FooterText:    d.Get("footer_text").(string),
		FooterLink:    d.Get("footer_link").(string),
	}

	if d.Get("custom_domain_name").(string) != "" {
		webformUpdateReq.IsCname = true
	}

	memailon := d.Get("email_on").([]interface{})
	emailon := make([]string, len(memailon))
	for i, v := range memailon {
		emailon[i] = v.(string)
	}
	webformUpdateReq.EmailOn = emailon

	mservices := d.Get("services").([]interface{})

	var services []api.WFService
	err := Decode(mservices, &services)
	if err != nil {
		return diag.FromErr(err)
	}
	webformUpdateReq.Services = services

	mseverity := d.Get("severity").([]interface{})
	var severity []api.WFSeverity
	err = Decode(mseverity, &severity)
	if err != nil {
		return diag.FromErr(err)
	}
	webformUpdateReq.Severity = severity

	mtags := d.Get("tags").(map[string]interface{})
	tags := make(map[string]string, len(*&mtags))
	for k, v := range *&mtags {
		tags[k] = v.(string)
	}

	webformUpdateReq.Tags = tags

	_, err = client.UpdateWebform(ctx, d.Get("team_id").(string), d.Id(), &webformUpdateReq)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceWebformRead(ctx, d, meta)
}

func resourceWebformDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	teamID, ok := d.GetOk("team_id")
	if !ok {
		return diag.Errorf("invalid team id provided")
	}
	_, err := client.DeleteWebform(ctx, teamID.(string), d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
