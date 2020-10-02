---
layout: "squadcast"
page_title: "squadcast: squadcast_escalation_policy"
sidebar_current: "docs-squadcast-datasource-escalation-policy"
description: |-
  Provides information about a Escalation Policy.

  This data source can be helpful when an escalation policy is handled outside(for eg: created on webapp/api's) Terraform but you still want to reference it in other resources.
---

# squadcast\_escalation_policy

Use this data source to get information about a specific escalation policy that you can use for other squadcast resources.

## Example Usage

```hcl
data "squadcast_escalation_policy" "test" {
  name = "Engineering Escalation Policy"
}

resource "squadcast_service" "test" {
  name                    = "My Web App"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = data.squadcast_escalation_policy.test.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to use to find an escalation policy in the squadcast API.

## Attributes Reference
* `id` - The ID of the found escalation policy.
* `name` - The short name of the found escalation policy.

[1]: https://v2.developer.squadcast.com/v2/page/api-reference#!/Escalation_Policies/get_escalation_policies
