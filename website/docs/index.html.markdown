---
layout: "squadcast"
page_title: "Provider: squadcast"
sidebar_current: "docs-squadcast-index"
description: |-
  Squadcast is an incident management and response tool.
---

# Squadcast Provider

Squadcast is an end-to-end incident management software that's designed to help you promote SRE best practices.

The provider configuration block accepts the following argument:

* ``squadcast_token`` - (Required) Refresh token of your Squadcast profile. 

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl

terraform {
  required_providers {
    squadcast = {
      source  = "SquadcastHub/squadcast"
    }
  }
}

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
  alert_source = "datadog"
}
```

!> **Warning:** Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system.

Token can also be passed using Environment variable
```sh
export squadcast_token=YOUR_TOKEN_HERE
```