data "squadcast_team" "example_resource_name" {
  name = "example team name"
}

data "squadcast_user" "example_resource_name" {
  email = "test@example.com"
}

data "squadcast_service" "example_resource_name" {
  name = "example service name"
}

data "squadcast_escalation_policy" "example_resource_name" {
  name = "example escalation policy name"
}

data "squadcast_squad" "example_resource_name" {
  name = "example squad name"
}

resource "squadcast_routing_rules" "example_resource_name" {
  team_id    = data.squadcast_team.example_resource_name.id
  service_id = data.squadcast_service.example_resource_name.id

  rules {
    is_basic   = false
    expression = "payload[\"event_id\"] == 40"

    route_to_id   = data.squadcast_user.example_resource_name.id / data.squadcast_squad.example_resource_name.id / data.squadcast_escalation_policy.example_resource_name.id
    route_to_type = "user/squad/escalation_policy"
  }

  rules {
    is_basic = true

    basic_expressions {
      lhs = "payload[\"foo\"]"
      rhs = "bar"
    }

    route_to_id   = data.squadcast_user.example_resource_name.id / data.squadcast_squad.example_resource_name.id / data.squadcast_escalation_policy.example_resource_name.id
    route_to_type = "user/squad/escalation_policy"
  }
}