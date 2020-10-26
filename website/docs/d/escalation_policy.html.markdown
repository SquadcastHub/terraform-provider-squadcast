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

resource "squadcast_service" "roz" {
  name = "datadog_service11"
  description = "Integrating Datadog with Squadcast....."
  escalation_policy_id =  data.squadcast_escalation_policy.rozd.id  
  email_prefix = "xya"
  alert_source = "datadog"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to use to find an escalation policy in the squadcast API.

## Attributes Reference
* `id` - The ID of the found escalation policy.

[1]: https://v2.developer.squadcast.com/v2/page/api-reference#!/Escalation_Policies/get_escalation_policies
