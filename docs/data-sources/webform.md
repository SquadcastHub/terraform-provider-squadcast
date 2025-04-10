---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "squadcast_webform Data Source - terraform-provider-squadcast"
subcategory: ""
description: |-
  Squadcast Webforms https://support.squadcast.com/webforms/webforms allows organizations to expand their customer support by hosting public Webforms, so their customers can quickly create an alert from outside the Squadcast ecosystem. Not only this, but internal stakeholders can also leverage Webforms for easy alert creation. Use this data source to get information about a specific webform.
---

# squadcast_webform (Data Source)

[Squadcast Webforms](https://support.squadcast.com/webforms/webforms) allows organizations to expand their customer support by hosting public Webforms, so their customers can quickly create an alert from outside the Squadcast ecosystem. Not only this, but internal stakeholders can also leverage Webforms for easy alert creation. Use this data source to get information about a specific webform.

## Example Usage

```terraform
data "squadcast_webform" "webform" {
  name    = "webformName"
  team_id = "team id"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the Webform.
- `team_id` (String) Team id.

### Read-Only

- `custom_domain_name` (String) Custom domain name (URL).
- `description` (String) Description of the Webform.
- `email_on` (List of String) Defines when to send email to the reporter (triggered, acknowledged, resolved).
- `footer_link` (String) Footer link.
- `footer_text` (String) Footer text.
- `header` (String) Webform header.
- `id` (Number) Webform id.
- `input_field` (List of Object) Input Fields added to Webforms. Added as tags to incident based on selection. (see [below for nested schema](#nestedatt--input_field))
- `owner` (List of Object) Form owner. (see [below for nested schema](#nestedatt--owner))
- `public_url` (String) Public URL of the Webform.
- `services` (List of Object) Services added to Webform. (see [below for nested schema](#nestedatt--services))
- `tags` (Map of String) Webform Tags.
- `title` (String) Webform title (public).

<a id="nestedatt--input_field"></a>
### Nested Schema for `input_field`

Read-Only:

- `label` (String)
- `options` (List of String)


<a id="nestedatt--owner"></a>
### Nested Schema for `owner`

Read-Only:

- `id` (String)
- `name` (String)
- `type` (String)


<a id="nestedatt--services"></a>
### Nested Schema for `services`

Read-Only:

- `alias` (String)
- `name` (String)
- `service_id` (String)
