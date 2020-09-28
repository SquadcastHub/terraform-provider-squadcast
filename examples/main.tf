
terraform {
  required_providers {
    squadcast = {
      versions = ["0.4"]
      source = "squadcast.com/tp/squadcast"
    }
  }
}


provider "squadcast" {
  # squadcast_token = "2287c30509d3c976e0a398f89325392e959dcb1c02432e1ca80987927832830034e442b3cce9c1e118fed3f897439c877cb253cc7dd89d085b9c8fd15a8fe8d1"
}

data "squadcast_escalation_policy" "rozd" {
  name = "example"
}

resource "squadcast_service" "roz" {
  name = "datadog_service11"
  description = "Integrating Datadog with Squadcast10"
  escalation_policy_id =  data.squadcast_escalation_policy.rozd.id  
  email_prefix = "xya10@gmal.com"
}