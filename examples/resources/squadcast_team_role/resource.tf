data "squadcast_team" "example_resource_name" {
  name = "example team name"
}

resource "squadcast_team_role" "test" {
  name      = "test"
  team_id   = data.squadcast_team.example.id
  abilities = ["create-escalation-policies", "read-escalation-policies", "update-escalation-policies"]
}