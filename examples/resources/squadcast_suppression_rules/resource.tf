data "squadcast_team" "example_team" {
  name = "exammple team name"
}

data "squadcast_service" "example_service" {
  name = "example service name"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_suppression_rules" "example_suppression_rules" {
  team_id    = data.squadcast_team.example_team.id
  service_id = data.squadcast_service.example_service.id

  rules {
    is_basic    = false
    description = "not basic"
    expression  = "payload[\"event_id\"] == 40"
  }
}
