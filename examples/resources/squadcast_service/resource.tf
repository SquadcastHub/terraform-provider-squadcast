data "squadcast_team" "example_team_resource" {
  name = "example team name"
}

data "squadcast_escalation_policy" "example_escalaion_policy_resource" {
  name = "example escalation policy name"
  team_id = data.squadcast_team.example_team_resource.id
}
resource "squadcast_service" "example_service_resource" {
  name                 = "example service name"
  team_id              = data.squadcast_team.example_team_resource.id
  escalation_policy_id = data.squadcast_escalation_policy.example_escalaion_policy_resource.id
  email_prefix         = "example-service-email"
}