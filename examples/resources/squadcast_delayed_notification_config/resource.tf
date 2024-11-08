data "squadcast_user" "example_user" {
  email = "user@example.com"
}

data "squadcast_team" "example_team" {
  name = "example team name"
}

data "squadcast_service" "example_service" {
  name    = "example service name"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_delayed_notification_config" "fixed_timeslot" {
  service_id = data.squadcast_service.example_service.id
  is_enabled = true
  timezone   = "Asia/Kolkata"

  fixed_timeslot_config {
    start_time  = "09:00"
    end_time    = "18:00"
    repeat_days = ["sunday", "monday", "tuesday"]
  }
  assigned_to {
    id   = data.squadcast_user.user.id
    type = "user"
  }
}

resource "squadcast_delayed_notification_config" "custom_timeslot" {
  service_id               = data.squadcast_service.example_service.id
  is_enabled               = true
  timezone                 = "Asia/Kolkata"
  custom_timeslots_enabled = true
  custom_timeslots {
    day_of_week = "sunday"
    start_time  = "10:15"
    end_time    = "20:00"
  }
  custom_timeslots {
    day_of_week = "monday"
    start_time  = "13:15"
    end_time    = "23:59"
  }
  custom_timeslots {
    day_of_week = "tuesday"
    start_time  = "12:15"
    end_time    = "20:59"
  }
  custom_timeslots {
    day_of_week = "wednesday"
    start_time  = "10:15"
    end_time    = "23:59"
  }


  assigned_to {
    id   = data.squadcast_user.example_user.id
    type = "user"
  }
}
