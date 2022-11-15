data "squadcast_team" "example_team_resource" {
  name = "example team name"
}

resource "squadcast_schedule" "example_schedule_resource" {
  name    = "example schedule name"
  team_id = data.squadcast_team.example_team_resource.id
  color   = "#9900ef"
}