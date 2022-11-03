package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type Contact struct {
	DialCode    string `json:"dial_code" tf:"-"`
	PhoneNumber string `json:"phone_number" tf:"-"`
}

type Ability struct {
	ID string `json:"id" tf:"-"`
	// Name    string `json:"name" tf:"-"`
	Slug string `json:"slug" tf:"-"`
	// Default bool   `json:"default" tf:"-"`
}

type PersonalNotificationRule struct {
	Type         string `json:"type" tf:"type"`
	DelayMinutes int    `json:"time" tf:"delay_minutes"`
}

func (r *PersonalNotificationRule) Encode() (tf.M, error) {
	return tf.Encode(r)
}

type OncallReminderRule struct {
	Type         string `json:"type" tf:"type"`
	DelayMinutes int    `json:"time" tf:"delay_minutes"`
}

func (r *OncallReminderRule) Encode() (tf.M, error) {
	return tf.Encode(r)
}

type DataSourceUser struct {
	AbilitiesSlugs            []string                    `json:"-" tf:"abilities"`
	Name                      string                      `json:"-" tf:"name"`
	PhoneNumber               string                      `json:"-" tf:"phone"`
	ID                        string                      `json:"id" tf:"id"`
	Abilities                 []*Ability                  `json:"abilities" tf:"-"`
	Bio                       string                      `json:"bio" tf:"-"`
	Contact                   Contact                     `json:"contact" tf:"-"`
	Email                     string                      `json:"email" tf:"email"`
	FirstName                 string                      `json:"first_name" tf:"first_name"`
	IsEmailVerified           bool                        `json:"email_verified" tf:"is_email_verified"`
	IsInGracePeriod           bool                        `json:"in_grace_period" tf:"-"`
	IsOverrideDnDEnabled      bool                        `json:"is_override_dnd_enabled" tf:"is_override_dnd_enabled"`
	IsPhoneVerified           bool                        `json:"phone_verified" tf:"is_phone_verified"`
	IsTrialSignup             bool                        `json:"is_trial_signup" tf:"-"`
	LastName                  string                      `json:"last_name" tf:"last_name"`
	OncallReminderRules       []*OncallReminderRule       `json:"oncall_reminder_rules" tf:"-"`
	PersonalNotificationRules []*PersonalNotificationRule `json:"notification_rules" tf:"-"`
	Role                      string                      `json:"role" tf:"role"`
	TimeZone                  string                      `json:"time_zone" tf:"time_zone"`
	Title                     string                      `json:"title" tf:"-"`
}

func (u *DataSourceUser) Encode() (tf.M, error) {
	u.Name = u.FirstName + " " + u.LastName

	if u.Contact.DialCode != "" && u.Contact.PhoneNumber != "" {
		u.PhoneNumber = u.Contact.DialCode + u.Contact.PhoneNumber
	}

	for _, v := range u.Abilities {
		u.AbilitiesSlugs = append(u.AbilitiesSlugs, v.Slug)
	}

	m, err := tf.Encode(u)
	if err != nil {
		return nil, err
	}

	sort.Strings(u.AbilitiesSlugs)
	m["abilities"] = u.AbilitiesSlugs

	rules, err := tf.EncodeSlice(u.OncallReminderRules)
	if err != nil {
		return nil, err
	}
	m["oncall_reminder_rules"] = rules

	rules, err = tf.EncodeSlice(u.PersonalNotificationRules)
	if err != nil {
		return nil, err
	}
	m["notification_rules"] = rules

	return m, nil
}

type ResourceUser struct {
	ID        string `json:"id" tf:"id"`
	Email     string `json:"email" tf:"email"`
	FirstName string `json:"first_name" tf:"first_name"`
	LastName  string `json:"last_name" tf:"last_name"`
	Role      string `json:"role" tf:"role"`

	Abilities      []*Ability `json:"abilities" tf:"-"`
	AbilitiesSlugs []string   `json:"-" tf:"abilities"`
}

func (u *ResourceUser) Encode() (tf.M, error) {
	for _, v := range u.Abilities {
		u.AbilitiesSlugs = append(u.AbilitiesSlugs, v.Slug)
	}

	m, err := tf.Encode(u)
	if err != nil {
		return nil, err
	}

	sort.Strings(u.AbilitiesSlugs)
	m["abilities"] = u.AbilitiesSlugs

	return m, nil
}

func (client *Client) GetUserById(ctx context.Context, id string) (*ResourceUser, error) {
	url := fmt.Sprintf("%s/users/%s", client.BaseURLV3, id)

	return Request[any, ResourceUser](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) GetUserByEmail(ctx context.Context, email string) (*DataSourceUser, error) {
	url := fmt.Sprintf("%s/users?email=%s", client.BaseURLV3, url.QueryEscape(email))

	return Request[any, DataSourceUser](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) ListUsers(ctx context.Context) ([]*ResourceUser, error) {
	url := fmt.Sprintf("%s/users", client.BaseURLV3)

	return RequestSlice[any, ResourceUser](http.MethodGet, url, client, ctx, nil)
}

type CreateUserReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

type UpdateUserReq struct {
	Role string `json:"role"`
}

type CreateUpdateUserResp struct {
	ID string `json:"id" tf:"id"`
}

type UpdateUserAbilitiesReq struct {
	UserID    string   `json:"user_id"`
	Abilities []string `json:"abilities"`
}

func (client *Client) CreateUser(ctx context.Context, req *CreateUserReq) (*CreateUpdateUserResp, error) {
	url := fmt.Sprintf("%s/users", client.BaseURLV3)

	return Request[CreateUserReq, CreateUpdateUserResp](http.MethodPost, url, client, ctx, req)
}

func (client *Client) UpdateUser(ctx context.Context, id string, req *UpdateUserReq) (*CreateUpdateUserResp, error) {
	url := fmt.Sprintf("%s/users/%s", client.BaseURLV3, id)

	return Request[UpdateUserReq, CreateUpdateUserResp](http.MethodPut, url, client, ctx, req)
}

func (client *Client) DeleteUser(ctx context.Context, id string) (*any, error) {
	url := fmt.Sprintf("%s/users/%s", client.BaseURLV3, id)

	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}

func (client *Client) UpdateUserAbilities(ctx context.Context, req *UpdateUserAbilitiesReq) (*any, error) {
	url := fmt.Sprintf("%s/users/abilities", client.BaseURLV3)

	type wrapped struct {
		Data []*UpdateUserAbilitiesReq `json:"data"`
	}
	bulkReq := wrapped{[]*UpdateUserAbilitiesReq{req}}

	return Request[wrapped, any](http.MethodPut, url, client, ctx, &bulkReq)
}
