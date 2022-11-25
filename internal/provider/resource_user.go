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

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "User resource.",

		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceUserImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "User id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"first_name": {
				Description:  "User first name.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"last_name": {
				Description:  "User last name.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"email": {
				Description:  "User email.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				ForceNew:     true,
			},
			"role": {
				Description:  "User role.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"user", "stakeholder", "account_owner"}, false),
			},
			"abilities": {
				Description: "user abilities/permissions.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceUserImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*api.Client)
	email := d.Id()

	user, err := client.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	d.SetId(user.ID)

	return []*schema.ResourceData{d}, nil
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	role := d.Get("role").(string)
	abilities := tf.ListToSlice[string](d.Get("abilities"))

	if role == "stakeholder" && len(abilities) != 0 {
		return diag.Errorf("stakeholders cannot have special abilities")
	}

	tflog.Info(ctx, "Creating user", tf.M{})
	user, err := client.CreateUser(ctx, &api.CreateUserReq{
		FirstName: d.Get("first_name").(string),
		LastName:  d.Get("last_name").(string),
		Email:     d.Get("email").(string),
		Role:      role,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(user.ID)

	if d.HasChange("abilities") {
		_, err := client.UpdateUserAbilities(ctx, &api.UpdateUserAbilitiesReq{
			UserID:    user.ID,
			Abilities: abilities,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	id := d.Id()

	tflog.Info(ctx, "Reading user", tf.M{
		"id": id,
	})
	user, err := client.GetUserById(ctx, id)
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = tf.EncodeAndSet(user, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	if d.HasChangesExcept("role", "abilities") {
		diag.Errorf("cannot change any attribute other than `role` or `abilities` for user `%s`. They can be only modified by the respective user in their profile page.", d.Get("email").(string))
	}

	role := d.Get("role").(string)
	abilities := tf.ListToSlice[string](d.Get("abilities"))

	if role == "stakeholder" && len(abilities) != 0 {
		return diag.Errorf("stakeholders cannot have special abilities")
	}

	if d.HasChange("role") {
		_, err := client.UpdateUser(ctx, d.Id(), &api.UpdateUserReq{
			Role: role,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("abilities") {
		_, err := client.UpdateUserAbilities(ctx, &api.UpdateUserAbilitiesReq{
			UserID:    d.Id(),
			Abilities: abilities,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceUserRead(ctx, d, meta)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.Client)

	_, err := client.DeleteUser(ctx, d.Id())
	if err != nil {
		if api.IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
