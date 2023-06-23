data "squadcast_team" "example_team" {
  name = "example team name"
}
data "squadcast_user" "example_user" {
  email = "test@example.com"
}
resource "squadcast_schedule_v2" "schedule_test" {
  name = "test schedule"
  description =  "test schedule"
  timezone = "Asia/Kolkata"
  team_id = data.squadcast_team.example_team.id
  entity_owner {
    id = data.squadcast_user.example_user.id
    type = "user"
  }
  tags {
    key = "testkey"
    value = "testval"
  }
  tags {
    key = "testkey2"
    value = "testval2"
  }
}
