resource "squadcast_suppression_rules" "test" {
  team_id    = "owner_id"
  service_id = "service_id"

  rules {
    is_basic    = false
    description = "not basic"
    expression  = "payload[\"event_id\"] == 40"
  }
}