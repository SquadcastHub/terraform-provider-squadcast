---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "squadcast_slo Resource - terraform-provider-squadcast"
subcategory: ""
description: |-
  squadcast_slo manages an SLO.
---

# squadcast_slo (Resource)

`squadcast_slo` manages an SLO.

## Example Usage

```terraform
resource "squadcast_slo" "test" {
  name               = "checkout-flow"
  description        = "Slo for checkout flow"
  target_slo         = 99.99
  service_ids        = ["615d3e23aff6885f46d291be"]
  slis               = ["latency", "high-err-rate"]
  time_interval_type = "rolling"
  duration_in_days   = 7
  org_id             = "604592dabc35ea0008bb0584"

  rules {
    name = "breached_error_budget"
  }

  rules {
    name      = "remaining_error_budget"
    threshold = 11
  }

  rules {
    name      = "unhealthy_slo"
    threshold = 1
  }

  notify {
    users = ["5e1c2309342445001180f9c2", "617793e650d38001057faaaf"]
  }

  owner_type = "team"
  team_id    = "611262fcd5b4ea846b534a8a"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the SLO.
- `service_ids` (List of String) Service IDs associated with the SLO.Only incidents from the associated services can be promoted as SLO violating incident
- `slis` (List of String) List of indentified SLIs for the SLO
- `target_slo` (Number) The target SLO for the time period.
- `team_id` (String) The team which SLO resource belongs to
- `time_interval_type` (String) Type of the SLO. Values can either be "rolling" or "fixed"

### Optional

- `description` (String) Description of the SLO.
- `duration_in_days` (Number) Tracks SLO for the last x days. Required only when SLO time interval type set to "rolling"
- `end_time` (String) SLO end time. Required only when SLO time interval type set to "fixed"
- `notify` (Block List, Max: 1) Notification rules for SLO violationUser can either choose to create an incident or get alerted via email (see [below for nested schema](#nestedblock--notify))
- `rules` (Block List) SLO monitoring checks has rules for monitoring any SLO violation(Or warning signs) (see [below for nested schema](#nestedblock--rules))
- `start_time` (String) SLO start time. Required only when SLO time interval type set to "fixed"

### Read-Only

- `id` (String) The ID of the SLO.

<a id="nestedblock--notify"></a>
### Nested Schema for `notify`

Optional:

- `service_id` (String) The ID of the service in which the user want to create an incident
- `squad_ids` (List of String) List of Squad ID's who should be alerted via email.
- `user_ids` (List of String) List of user ID's who should be alerted via email.

Read-Only:

- `id` (Number) The ID of the notification rule
- `slo_id` (Number) The ID of the SLO.


<a id="nestedblock--rules"></a>
### Nested Schema for `rules`

Required:

- `name` (String) The name of monitoring check."Supported values are "breached_error_budget", "unhealthy_slo","increased_false_positives", "remaining_error_budget"

Optional:

- `threshold` (Number) Threshold for the monitoring checkOnly supported for rules name "increased_false_positives" and "remaining_error_budget"

Read-Only:

- `id` (Number) The ID of the monitoring rule
- `is_checked` (Boolean) Is checked?
- `slo_id` (Number) The ID of the SLO

