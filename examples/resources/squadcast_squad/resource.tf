resource "squadcast_squad" "test" {
  name       = "test squad"
  team_id    = "owner_id"
  member_ids = ["user_id"]
}