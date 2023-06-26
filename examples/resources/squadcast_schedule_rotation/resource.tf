data "squadcast_team" "example_team" {
  name = "example team name"
}
data "squadcast_user" "example_user" {
  email = "test@example.com"
}
data "squadcast_user" "example_user_2" {
  email = "test2@example.com"
}

data "squadcast_schedule_v2" "get_schedule" {
  name = "Test Schedule"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_schedule_rotation" "rotations_with_daily_period" {
    schedule_id = data.squadcast_schedule_v2.get_schedule.id
    name = "Test Rotation"
    start_date = "2023-07-01T00:00:00Z"
    period = "daily"
    shift_timeslots {
        start_hour = 10
        start_minute = 30
        duration = 720
    }
    change_participants_frequency = 1
    change_participants_unit = "rotation"
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
    ends_after_iterations = 2
}

resource "squadcast_schedule_rotation" "rotations_with_custom_period" {
    schedule_id = data.squadcast_schedule_v2.get_schedule.id
    name = "Test Rotation 2"
    start_date = "2023-06-13T00:00:00Z"
    period = "custom"
    shift_timeslots {
        start_hour = 10
        start_minute = 0
        duration = 1440
        day_of_week = "saturday"
    }
    shift_timeslots {
        start_hour = 12
        start_minute = 30
        duration = 720
        day_of_week = "sunday"
    }
    change_participants_frequency = 1
    change_participants_unit = "rotation"
    custom_period_frequency = 1
    custom_period_unit = "week"
    participant_groups {
        participants {
            id = data.squadcast_user.example_user.id
            type = "user"
        }
    }
    participant_groups {
        participants {
            id = data.squadcast_user.example_user_2.id
            type = "user"
        }
    }
    end_date =  "2023-08-31T00:00:00Z"
}
