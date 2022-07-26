data "squadcast_team" "example_resource_name" {
  name = "example team name"
}

resource "squadcast_schedule" "example_resource_name" {
  name    = "example schedule name"
  team_id = data.squadcast_team.example_resource_name.id
  color   = "#9900ef"
}