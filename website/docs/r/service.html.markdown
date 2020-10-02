---
layout: "squadcast"
page_title: "squadcast: squadcast_service"
sidebar_current: "docs-squadcast_service"
description: |-
  The squadcast_service resource allows squadcast services to be managed by Terraform.
---

# squadcast\_service

[Services](https://support.squadcast.com/docs/adding-a-service-1) are at the core of Squadcast. A service represents an application or component that is crucial for your product or service. Services are created with an alert source integration through which incidents are triggered. Squadcast provides a Webhook URL to integrate with the tools you use.

## Example Usage

```hcl
data "squadcast_escalation_policy" "main" {
  name = "example"
}

resource "squadcast_service" "main" {
  name = "datadog_service"
  description = "Integrating Datadog with Squadcast"
  escalation_policy_id =  data.squadcast_escalation_policy.main.id
  email_prefix = "xyz"
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) This is the name of the test and the website to be monitored.

## Attributes Reference

The following attribute is exported:

* `id` - A unique identifier for the escalation policy
