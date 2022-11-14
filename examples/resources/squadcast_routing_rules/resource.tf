data "squadcast_team" "example_team_resource" {
  name = "example team name"
}

data "squadcast_user" "example_user_resource" {
  email = "test@example.com"
}

data "squadcast_service" "example_service_resource" {
  name = "example service name"
  team_id = data.squadcast_team.example_team_resource.id
}

data "squadcast_escalation_policy" "example_escalaion_policy_resource" {
  name = "example escalation policy name"
  team_id = data.squadcast_team.example_team_resource.id
}

data "squadcast_squad" "example_squad_resource" {
  name = "example squad name"
  team_id = data.squadcast_team.example_team_resource.id
}

resource "squadcast_routing_rules" "example_routing_rules_resource" {
  team_id    = data.squadcast_team.example_team_resource.id
  service_id = data.squadcast_service.example_service_resource.id

  rules {
    is_basic   = false
    expression = "payload[\"event_id\"] == 40"

    route_to_id   = data.squadcast_user.example_user_resource.id / data.squadcast_squad.example_squad_resource.id / data.squadcast_escalation_policy.example_escalaion_policy_resource.id
    route_to_type = "user/squad/escalation_policy"
  }

  rules {
    is_basic = true

    basic_expressions {
      lhs = "payload[\"foo\"]"
      rhs = "bar"
    }

    route_to_id   = data.squadcast_user.example_user_resource.id / data.squadcast_squad.example_squad_resource.id / data.squadcast_escalation_policy.example_escalaion_policy_resource.id
    route_to_type = "user/squad/escalation_policy"
  }
}