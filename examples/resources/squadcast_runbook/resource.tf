data "squadcast_team" "example_resource_name" {
  name = "example team name"
}

resource "squadcast_runbook" "example_resource_name" {
  name    = "example runbook name"
  team_id = data.squadcast_team.example_resource_name.id

  steps {
    content = "some text here"
  }

  steps {
    content = "some text here 2"
  }
}