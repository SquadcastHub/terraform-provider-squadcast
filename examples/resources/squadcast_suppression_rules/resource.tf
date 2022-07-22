data "squadcast_team" "example" {
  name = "test"
}

data "squadcast_service" "example" {
  name = "test-parent"
}

resource "squadcast_suppression_rules" "test" {
  team_id    = data.squadcast_team.example.id
  service_id = data.squadcast_service.example.id

  rules {
    is_basic    = false
    description = "not basic"
    expression  = "payload[\"event_id\"] == 40"
  }
}