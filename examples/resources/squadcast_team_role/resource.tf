data "squadcast_team" "example_team_resource" {
  name = "example team name"
}

resource "squadcast_team_role" "example_team_role_resource" {
  name      = "test"
  team_id   = data.squadcast_team.example_team_resource.id
  abilities = ["create-escalation-policies", "read-escalation-policies", "update-escalation-policies"]
}