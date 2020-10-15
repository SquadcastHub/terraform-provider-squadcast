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
data "squadcast_escalation_policy" "rozd" {
  name = "example"
}

resource "squadcast_service" "roz" {
  name = "datadog_service11"
  description = "Integrating Datadog with Squadcast....."
  escalation_policy_id =  data.squadcast_escalation_policy.rozd.id  
  email_prefix = "xya10@gmal.com"
  alert_source = "datadog"
}

output "webhook_url" {
  value = squadcast_service.roz.webhook_url
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) This is the name of the test and the website to be monitored.
* `description` - (Optional) Short description of the service.
* `escalation_policy_id` - Object id of the service
* `email_prefix` - Email prefix for the service 
* `alert_source` (Required) The name of the alert source being used with the service. eg: `datadog`, `pingdom`

## Attributes Reference

The following attribute is exported:

* `id` - A unique identifier for the escalation policy
* `webhook_url` - Webhook url for the alert source integration.
