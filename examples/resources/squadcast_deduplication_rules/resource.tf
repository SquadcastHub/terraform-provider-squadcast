data "squadcast_team" "example_team" {
  name = "example team name"
}

data "squadcast_service" "example_service" {
  name = "example service name"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_deduplication_rules" "example_deduplication_rules" {
  team_id    = data.squadcast_team.example_team.id
  service_id = data.squadcast_service.example_service.id

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
