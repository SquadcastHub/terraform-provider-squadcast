data "squadcast_team" "example" {
  name = "test"
}

data "squadcast_escalation_policy" "example" {
  name = "test"
}
resource "squadcast_service" "test_parent" {
  name                 = "test-service-parent"
  team_id              = data.squadcast_team.example.id
  escalation_policy_id = data.squadcast_escalation_policy.example.id
  email_prefix         = "test-service-parent"
}