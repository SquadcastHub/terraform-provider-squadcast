data "squadcast_user" "example_user" {
  email = "user@example.com"
}

resource "squadcast_team" "example_team" {
  name            = "example team name"
  default_user_id = data.squadcast_user.example_user.id
}
