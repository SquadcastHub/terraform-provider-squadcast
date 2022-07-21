resource "squadcast_runbook" "test" {
  name    = "test runbook"
  team_id = "owner_id"

  steps {
    content = "some text here"
  }

  steps {
    content = "some text here 2"
  }
}