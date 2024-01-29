data "squadcast_user" "example_user" {
  email = "test@example.com"
}

data "squadcast_team" "example_team" {
  name = "example team name"
}

data "squadcast_escalation_policy" "example_escalaion_policy" {
  name    = "example escalation policy name"
  team_id = data.squadcast_team.example_team.id
}
resource "squadcast_service" "example_service" {
  name                 = "example service name"
  team_id              = data.squadcast_team.example_team.id
  escalation_policy_id = data.squadcast_escalation_policy.example_escalaion_policy.id
  email_prefix         = "example-service-email"
  maintainer {
    id   = data.squadcast_user.example_user.id
    type = "user"
  }
  tags {
    key   = "testkey"
    value = "testval"
  }
  tags {
    key   = "testkey2"
    value = "testval2"
  }
  alert_sources = ["example-alert-source"]
  slack_channel_id = "D0KAQDEPSH"
}

resource "squadcast_service" "example_service_with_delay_notification" {
  name                 = "Fixed timeslot service"
  team_id              = data.squadcast_team.example_team.id
  email_prefix         = "example-service-email"
  escalation_policy_id = data.squadcast_escalation_policy.example_escalaion_policy.id
  description          = "example service description"

  delay_notification_config {
    is_enabled = true
    timezone   = "Asia/Kolkata"
  
    fixed_timeslot_config {
      start_time   = "09:00"
      end_time     = "18:00"
      repeat_days  = ["sunday","monday","tuesday"]
    }
    assigned_to {
      type = "default_escalation_policy"
    }
  }
}

resource "squadcast_service" "example_service_with_custom_timeslot" {
  name                 = "Custom timeslot service"
  team_id              = data.squadcast_team.example_team.id
  email_prefix         = "example-service-email"
  escalation_policy_id = data.squadcast_escalation_policy.example_escalaion_policy.id
  description          = "example service description"

  delay_notification_config {
    is_enabled = true
    timezone                      = "Asia/Kolkata"
    custom_timeslots_enabled      = true
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
}