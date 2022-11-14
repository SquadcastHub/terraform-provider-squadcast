data "squadcast_team" "example_resource_name" {
  name = "example test name"
}

data "squadcast_user" "example_user_resource" {
  email = "test@example.com"
}
resource "squadcast_squad" "example_squad_resource" {
  name       = "example squad name"
  team_id    = data.squadcast_team.example_team_resource.id
  member_ids = [data.squadcast_user.example_user_resource.id]
}