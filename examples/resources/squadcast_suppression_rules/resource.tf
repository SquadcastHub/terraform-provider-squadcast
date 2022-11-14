data "squadcast_team" "example_team_resource" {
  name = "exammple team name"
}

data "squadcast_service" "example_service_resource" {
  name = "example service name"
  team_id = data.squadcast_team.example_team_resource.id
}

resource "squadcast_suppression_rules" "example_suppression_rules_resource" {
  team_id    = data.squadcast_team.example_team_resource.id
  service_id = data.squadcast_service.example_service_resource.id

  rules {
    is_basic    = false
    description = "not basic"
    expression  = "payload[\"event_id\"] == 40"
  }
}