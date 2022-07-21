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

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information about a specific user that you can use for other Squadcast resources.",

		ReadContext: dataSourceUserRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "User id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"first_name": {
				Description: "User first name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_name": {
				Description: "User last name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "User name, automatically computed from first name and last name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description:  "User email.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"is_email_verified": {
				Description: "Denotes if the user has verified their email or not.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"phone": {
				Description: "User phone number.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_phone_verified": {
				Description: "Denotes if the user has verified their phone number or not.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"is_override_dnd_enabled": {
				Description: "Deprecated, this can be ignored.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"role": {
				Description: "User role.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"time_zone": {
				Description: "User time_zone.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"abilities": {
				Description: "Denotes the Permissions / abilities of the user.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"notification_rules": {
				Description: "User Personal Notification Rules.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Personal notification rule type.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"delay_minutes": {
							Description: "notification rule delay_minutes, (to be deprecated).",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
			"oncall_reminder_rules": {
				Description: "User's personal on-call reminder notification rules.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "oncall reminder rule type.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"delay_minutes": {
							Description: "oncall reminder rule delay_minutes.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	email := d.Get("email").(string)

	tflog.Info(ctx, "Reading user", tf.M{
		"email": email,
	})
	user, err := client.GetUserByEmail(ctx, email)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(user, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
