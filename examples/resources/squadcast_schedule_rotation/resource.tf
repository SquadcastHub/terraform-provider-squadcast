data "squadcast_team" "example_team" {
  name = "example team name"
}
data "squadcast_user" "example_user" {
  email = "test@example.com"
}
data "squadcast_user" "example_user_2" {
  email = "test2@example.com"
}

resource "squadcast_schedule_v2" "schedule_test" {
  name = "Test schedule"
  description =  "test schedule"
  timezone = "Asia/Kolkata"
  team_id = data.squadcast_team.example_team.id
  entity_owner {
    id = data.squadcast_user.example_user.id
    type = "user"
  }
}

resource "squadcast_schedule_rotation" "rotations_with_custom_period" {
    schedule_id = squadcast_schedule_v2.schedule_test.id
    name = "Test Rotation"
    start_date = "2021-08-01T00:00:00Z"
    period = "custom"
    shift_timeslots {
        start_hour = 10
        start_minute = 0
        duration = 10
        day_of_week = "monday"
    }
    shift_timeslots {
        start_hour = 12
        start_minute = 30
        duration = 60
        day_of_week = "friday"
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
        participants {
            id = data.squadcast_user.example_user_2.id
            type = "user"
        }
    }
}

resource "squadcast_schedule_rotation" "rotations_with_daily_period" {
    schedule_id = squadcast_schedule_v2.schedule_test.id
    name = "Test Rotation 2"
    start_date = "2021-08-01T00:00:00Z"
    period = "daily"
    shift_timeslots {
        start_hour = 10
        start_minute = 30
        duration = 120
        day_of_week = "monday"
    }
    change_participants_frequency = 1
    change_participants_unit = "week"
    participant_groups {
        participants {
            id = data.squadcast_user.example_user.id
            type = "user"
        }
        participants {
            id = data.squadcast_user.example_user_2.id
            type = "user"
        }
    }
}