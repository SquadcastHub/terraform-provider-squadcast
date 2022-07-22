data "squadcast_team" "example" {
  name = "test"
}

data "squadcast_user" "example" {
  email = "test@example.com"
}
resource "squadcast_squad" "test" {
  name       = "test squad"
  team_id    = data.squadcast_team.example.id
  member_ids = [data.squadcast_user.example.id]
}