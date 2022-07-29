data "squadcast_team" "example_resource_name" {
  name = "example test name"
}

data "squadcast_user" "example_resource_name" {
  email = "test@example.com"
}
resource "squadcast_squad" "example_resource_name" {
  name       = "example squad name"
  team_id    = data.squadcast_team.example.id
  member_ids = [data.squadcast_user.example.id]
}