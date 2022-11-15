data "squadcast_team" "example_team_resource" {
  name = "example team name"
}

data "squadcast_service" "example_service_resource" {
  name = "example service name"
  team_id = data.squadcast_team.example_team_resource.id
}

resource "squadcast_deduplication_rules" "example_deduplication_rules_resource" {
  team_id    = data.squadcast_team.example_team_resource.id
  service_id = data.squadcast_service.example_service_resource.id

  rules {
    is_basic    = false
    description = "not basic"
    expression  = "payload[\"event_id\"] == 40"
  }

  rules {
    is_basic    = true
    description = "basic"

    basic_expressions {
      lhs = "payload[\"foo\"]"
      op  = "is"
      rhs = "bar"
    }
  }
}