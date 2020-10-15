
terraform {
  required_providers {
    squadcast = {
      versions = ["0.4"]
      source = "squadcast.com/tp/squadcast"
    }
  }
}


provider "squadcast" {
  squadcast_token = "YOUR_SQUADCAST_TOKEN"
}

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