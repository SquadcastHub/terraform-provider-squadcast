resource "squadcast_deduplication_rules" "test" {
  team_id    = "owner_id"
  service_id = "service_id"

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