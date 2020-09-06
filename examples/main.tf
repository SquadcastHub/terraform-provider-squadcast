
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

resource "squadcast_service" "rozd1" {
  name = "datadog_service"
  description = "Integrating Datadog with Squadcast"
  escalation_policy_id = "5f35a422ce4a1800086df873"
  email_prefix = "xyz@gmal.com"
}