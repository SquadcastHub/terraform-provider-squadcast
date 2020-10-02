---
layout: "squadcast"
page_title: "Provider: squadcast"
sidebar_current: "docs-squadcast-index"
description: |-
  Squadcast is an incident management and response tool.
---

# squadcast Provider

Squadcast is an end-to-end incident management software that's designed to help you promote SRE best practices.

The provider configuration block accepts the following argument:

* ``squadcast_token`` - (Required) Refresh token of your Squadcast profile.TThis can also be passed as a ``squadcast_token`` environment variable

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
provider "squadcast" {
  squadcast_token = "YOUR-SQUADCAST-TOKEN"
}

data "squadcast_escalation_policy" "main" {
  name = "example"
}

resource "squadcast_service" "main" {
  name = "datadog_service"
  description = "Integrating Datadog with Squadcast"
  escalation_policy_id =  data.squadcast_escalation_policy.main.id
  email_prefix = "xyz"
```
