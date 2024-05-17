data "squadcast_team" "example_team" {
  name = "Example Team"
}

resource "squadcast_global_oncall_reminder_rules" "example_config" {
  team_id = data.squadcast_team.example_team.id
  is_enabled = true

  rules {
    type = "Email"
    time = 60
  }
  rules {
    type = "Push"
    time = 60
  }
}