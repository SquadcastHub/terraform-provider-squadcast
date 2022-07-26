data "squadcast_team" "example_resource_name" {
  name = "example team name"
}

data "squadcast_escalation_policy" "example_resource_name" {
  name = "example escalation policy name"
}
resource "squadcast_service" "example_resource_name" {
  name                 = "example service name"
  team_id              = data.squadcast_team.example_resource_name.id
  escalation_policy_id = data.squadcast_escalation_policy.example_resource_name.id
  email_prefix         = "example-service-email"
}