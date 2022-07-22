data "squadcast_team" "example" {
  name = "test"
}

resource "squadcast_schedule" "test" {
  name    = "test schedule"
  team_id = data.squadcast_team.example.id
  color   = "#9900ef"
}