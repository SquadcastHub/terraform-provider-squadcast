data "squadcast_team" "example_team" {
  name = "example team name"
}

data "squadcast_user" "example_user" {
  email = "test@example.com"
}

data "squadcast_service" "example_service" {
  name = "example service name"
  team_id = data.squadcast_team.example_team.id
}

data "squadcast_escalation_policy" "example_escalaion_policy" {
  name = "example escalation policy name"
  team_id = data.squadcast_team.example_team.id
}

data "squadcast_squad" "example_squad" {
  name = "example squad name"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_routing_rules" "example_routing_rules" {
  team_id    = data.squadcast_team.example_team.id
  service_id = data.squadcast_service.example_service.id

  rules {
    is_basic   = false
    expression = "payload[\"event_id\"] == 40"

    route_to_id   = data.squadcast_user.example_user.id / data.squadcast_squad.example_squad.id / data.squadcast_escalation_policy.example_escalaion_policy.id
    route_to_type = "user/squad/escalation_policy"
  }

  rules {
    is_basic = true

    basic_expressions {
      lhs = "payload[\"foo\"]"
      rhs = "bar"
    }

    route_to_id   = data.squadcast_user.example_user.id / data.squadcast_squad.example_squad.id / data.squadcast_escalation_policy.example_escalaion_policy.id
    route_to_type = "user/squad/escalation_policy"
  }
}
