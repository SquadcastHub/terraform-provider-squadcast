resource "squadcast_team_member" "test" {
  team_id  = "owner_id"
  user_id  = "user_id"
  role_ids = ["role_id", "role_id"]
}