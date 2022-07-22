data "squadcast_team" "example" {
  name = "test"
}

data "squadcast_user" "example" {
  email = "test@example.com"
}

data "squadcast_service" "example" {
  name = "test-parent"
}

data "squadcast_escalation_policy" "example" {
  name = "test"
}

data "squadcast_squad" "example" {
  name = "test"
}

resource "squadcast_routing_rules" "test" {
  team_id    = data.squadcast_team.example.id
  service_id = data.squadcast_service.example.id

  rules {
    is_basic   = false
    expression = "payload[\"event_id\"] == 40"

    route_to_id   = data.squadcast_user.example.id / data.squadcast_squad.example.id / data.squadcast_escalation_policy.example.id
    route_to_type = "user/squad/escalation_policy"
  }

  rules {
    is_basic = true

    basic_expressions {
      lhs = "payload[\"foo\"]"
      rhs = "bar"
    }

    route_to_id   = data.squadcast_user.example.id / data.squadcast_squad.example.id / data.squadcast_escalation_policy.example.id
    route_to_type = "user/squad/escalation_policy"
  }
}