
terraform {
  required_providers {
    squadcast = {
      versions = ["0.4"]
      source = "squadcast.com/tp/squadcast"
    }
  }
}


provider "squadcast" {
  # squadcast_token = "xxx"
}

data "squadcast_escalation_policy" "roz" {
  name = "example"
}

resource "squadcast_service" "roz" {
  name = "datadog_service1"
  description = "Integrating Datadog with Squadcast"
  escalation_policy_id =  "data.squadcast_escalation_policy.roz.id" // "5f35a422ce4a1800086df873"
  email_prefix = "xya@gmal.com"
}