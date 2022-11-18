data "squadcast_team" "example_team" {
  name = "example team name"
}

resource "squadcast_team_role" "example_team_role" {
  name      = "test"
  team_id   = data.squadcast_team.example_team.id
  abilities = ["create-escalation-policies", "read-escalation-policies", "update-escalation-policies"]
}
