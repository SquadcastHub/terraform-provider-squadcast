data "squadcast_team" "example_team" {
  name = "example team name"
}

resource "squadcast_schedule" "example_schedule" {
  name    = "example schedule name"
  team_id = data.squadcast_team.example_team.id
  color   = "#9900ef"
}
