data "squadcast_user" "example_user" {
  email = "test@example.com"
}

data "squadcast_team" "example_team" {
  name = "example team name"
}

data "squadcast_escalation_policy" "example_escalaion_policy" {
  name = "example escalation policy name"
  team_id = data.squadcast_team.example_team.id
}
resource "squadcast_service" "example_service" {
  name                 = "example service name"
  team_id              = data.squadcast_team.example_team.id
  escalation_policy_id = data.squadcast_escalation_policy.example_escalaion_policy.id
  email_prefix          = "example-service-email"
  maintainer {
    id = data.squadcast_user.example_user.id
    type = "user"
  }
  tags {
    key = "testkey"
    value = "testval"
  }
  tags {
    key = "testkey2"
    value = "testval2"
  }
  alert_sources = ["example-alert-source"]
}
