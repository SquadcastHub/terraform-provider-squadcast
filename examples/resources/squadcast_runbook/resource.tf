data "squadcast_team" "example" {
  name = "test"
}

resource "squadcast_runbook" "test" {
  name    = "test runbook"
  team_id = data.squadcast_team.example.id

  steps {
    content = "some text here"
  }

  steps {
    content = "some text here 2"
  }
}