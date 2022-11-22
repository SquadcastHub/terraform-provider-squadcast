data "squadcast_team" "example_team" {
  name = "example team name"
}

data "squadcast_user" "example_user" {
  email = "test@example.com"
}

data "squadcast_team_role" "example_team_role" {
  name = "example role name"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_team_member" "example_team_member" {
  team_id  = data.squadcast_team.example_team.id
  user_id  = data.squadcast_user.example_user.id
  role_ids = [data.squadcast_team_role.example_team_role.id]
}
