data "squadcast_team" "example_team_resource" {
  name = "example team name"
}

data "squadcast_user" "example_user_resource" {
  email = "test@example.com"
}

data "squadcast_team_role" "example_team_role_resource" {
  name = "example role name"
  team_id = data.squadcast_team.example_team_resource.id
}

resource "squadcast_team_member" "example_team_member_resource" {
  team_id  = data.squadcast_team.example_team_resource.id
  user_id  = data.squadcast_user.example_user_resource.id
  role_ids = [data.squadcast_team_role.example_team_role_resource.id]
}