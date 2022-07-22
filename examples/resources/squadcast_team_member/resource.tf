data "squadcast_team" "example" {
  name = "test"
}

data "squadcast_user" "example" {
  email = "test@example.com"
}

data "squadcast_team_role" "example" {
  name = "test"
}

resource "squadcast_team_member" "test" {
  team_id  = data.squadcast_team.example.id
  user_id  = data.squadcast_user.example.id
  role_ids = [data.squadcast_team_role.example.id]
}