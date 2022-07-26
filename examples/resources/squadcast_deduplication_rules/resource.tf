data "squadcast_team" "example_resource_name" {
  name = "example team_name"
}

data "squadcast_service" "example_resource_name" {
  name = "example service name"
}

resource "squadcast_deduplication_rules" "example_resource_name" {
  team_id    = data.squadcast_team.exaexample_resource_namemple.id
  service_id = data.squadcast_service.example_resource_name.id

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