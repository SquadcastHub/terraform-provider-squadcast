data "squadcast_team" "example" {
  name = "test"
}

resource "squadcast_team_role" "test" {
  name      = "test"
  team_id   = data.squadcast_team.example.id
  abilities = []
}