resource "squadcast_routing_rules" "test" {
  team_id    = "team_id"
  service_id = "service_id"

  rules {
    is_basic   = false
    expression = "payload[\"event_id\"] == 40"

    route_to_id   = "user_id/squad_id/escalatio_policy_id"
    route_to_type = "user/squad/escalation_policy"
  }

  rules {
    is_basic = true

    basic_expressions {
      lhs = "payload[\"foo\"]"
      rhs = "bar"
    }

    route_to_id   = "user_id/squad_id/escalatio_policy_id"
    route_to_type = "user/squad/escalation_policy"
  }
}