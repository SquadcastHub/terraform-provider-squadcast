resource "squadcast_service" "test_parent" {
  name                 = "test-service-parent"
  team_id              = "owner_id"
  escalation_policy_id = "escalatio_policy_id"
  email_prefix         = "test-service-parent"
}