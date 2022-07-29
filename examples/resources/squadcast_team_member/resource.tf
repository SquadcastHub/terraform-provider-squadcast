data "squadcast_team" "example_resource_name" {
  name = "example team name"
}

data "squadcast_user" "example_resource_name" {
  email = "test@example.com"
}

data "squadcast_team_role" "example_resource_name" {
  name = "example role name"
}

resource "squadcast_team_member" "test" {
  team_id  = data.squadcast_team.example_resource_name.id
  user_id  = data.squadcast_user.example_resource_name.id
  role_ids = [data.squadcast_team_role.example_resource_name.id]
}