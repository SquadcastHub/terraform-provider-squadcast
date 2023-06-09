data "squadcast_team" "example_team" {
  name = "example team name"
}
data "squadcast_user" "example_user" {
  email = "test@example.com"
}

resource "squadcast_schedule_v2" "name" {
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
    color = "#ccc"
  }
  tags {
    key = "testkey2"
    value = "testval2"
    color = "green"
  }
  rotations {
    name = "Test Rotation"
    start_date = "2023-06-09T00:00:00Z"
    period = "custom"
    shift_timeslots {
        start_hour = 10
        start_minute = 0
        duration = 10
        day_of_week = "monday"
    }
    shift_timeslots {
        start_hour = 10
        start_minute = 0
        duration = 10
        day_of_week = "monday"
    }
    change_participants_frequency = 1
    change_participants_unit = "week"
    custom_period_frequency = 1
    custom_period_unit = "week"
    participant_groups {
        participants {
            id = data.squadcast_user.example_user.id
            type = "user"
        }
    }
  }
}