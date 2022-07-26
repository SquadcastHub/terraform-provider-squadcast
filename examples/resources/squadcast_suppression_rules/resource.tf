data "squadcast_team" "example_resource_name" {
  name = "exammple team name"
}

data "squadcast_service" "example_resource_name" {
  name = "example service name"
}

resource "squadcast_suppression_rules" "example_resource_name" {
  team_id    = data.squadcast_team.example_resource_name.id
  service_id = data.squadcast_service.example_resource_name.id

  rules {
    is_basic    = false
    description = "not basic"
    expression  = "payload[\"event_id\"] == 40"
  }
}