data "squadcast_team" "example_team_resource" {
  name = "example team name"
}

resource "squadcast_runbook" "example_runbook_resource" {
  name    = "example runbook name"
  team_id = data.squadcast_team.example_team_resource.id

  steps {
    content = "some text here"
  }

  steps {
    content = "some text here 2"
  }
}