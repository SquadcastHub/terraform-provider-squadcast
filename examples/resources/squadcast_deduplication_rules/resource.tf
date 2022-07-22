data "squadcast_team" "example" {
  name = "test"
}

data "squadcast_service" "example" {
  name = "test"
}

resource "squadcast_deduplication_rules" "test" {
  team_id    = data.squadcast_team.example.id
  service_id = data.squadcast_service.example.id

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