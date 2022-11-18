data "squadcast_team" "example_team" {
  name = "example test name"
}

data "squadcast_user" "example_user" {
  email = "test@example.com"
}
resource "squadcast_squad" "example_squad" {
  name       = "example squad name"
  team_id    = data.squadcast_team.example_team.id
  member_ids = [data.squadcast_user.example_user.id]
}
