data "squadcast_team" "example_team" {
  name = "example test name"
}

data "squadcast_user" "example_user" {
  email = "test@example.com"
}
data "squadcast_user" "example_user2" {
  email = "test2@example.com"
}

resource "squadcast_squad" "example_squad" {
  name       = "example squad name"
  team_id    = data.squadcast_team.example_team.id
  member_ids = [data.squadcast_user.example_user.id] # deprecated
}

# RBAC permission model
resource "squadcast_squad" "example_squad_rbac" {
  name       = "example rbac squad"
  team_id    = data.squadcast_team.example_team.id
  members {
      user_id = data.squadcast_user.example_user.id
  }
  members {
      user_id = data.squadcast_user.example_user_2.id
  }
}

# OBAC permission model
resource "squadcast_squad" "example_squad_obac" {
  name       = "example obac squad"
  team_id    = data.squadcast_team.example_team.id
  members {
      user_id = data.squadcast_user.example_user.id
      role = "owner"
  }
  members {
      user_id = data.squadcast_user.example_user_2.id
      role = "member"
  }
}

