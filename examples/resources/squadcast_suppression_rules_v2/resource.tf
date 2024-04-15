data "squadcast_team" "example_team" {
  name = "exammple team name"
}

data "squadcast_service" "example_service" {
  name    = "example service name"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_suppression_rules_v2" "example_basic_suppression_rules" {
  service_id = data.squadcast_service.example_service.id
  is_basic    = true
  description = "sample basic expression"
  basic_expressions {
    lhs = "abc"
    op = "is"
    rhs = "xyz"
  }
}

resource "squadcast_suppression_rules_v2" "example_suppression_rules" {
  service_id = data.squadcast_service.example_service.id
  is_basic    = false
  description = "not basic"
  expression  = "payload[\"event_id\"] == 40"
}

resource "squadcast_suppression_rules_v2" "example_time_based_suppression_rules" {
  service_id = data.squadcast_service.example_service.id
  is_basic    = false
  description = "not basic"
  expression  = "payload[\"event_id\"] == 40"
  timeslots {
    time_zone  = "Asia/Calcutta"
    start_time = "2022-04-08T06:22:14.975Z"
    end_time   = "2022-04-28T06:22:14.975Z"
    ends_on    = "2022-04-28T06:22:14.975Z"
    repetition = "none" # none, daily, weekly, monthly, custom
    is_allday  = false
    ends_never = true
  }
}


resource "squadcast_suppression_rules_v2" "example_time_based_suppression_rules_custom_repetition" {
  service_id = data.squadcast_service.example_service.id
  is_basic    = false
  description = "not basic"
  expression  = "payload[\"event_id\"] == 40"
  # custom repetition - daily
  timeslots {
    time_zone  = "Asia/Calcutta"
    start_time = "2022-04-08T06:22:14.975Z"
    end_time   = "2022-04-28T06:22:14.975Z"
    ends_on    = "2022-04-28T06:22:14.975Z"
    repetition = "custom"
    is_allday  = false
    ends_never = true
    custom {
      repeats       = "day"
      repeats_count = 2
    }
  }
  # custom repetition - weekly
  timeslots {
    time_zone  = "Asia/Calcutta"
    start_time = "2022-04-08T06:22:14.975Z"
    end_time   = "2022-04-28T06:22:14.975Z"
    ends_on    = "2022-04-28T06:22:14.975Z"
    repetition = "custom"
    is_allday  = false
    ends_never = true
    custom {
      repeats             = "week"
      repeats_count       = 4
      repeats_on_weekdays = [0, 1, 2, 3] # 0 - Sunday, 1 - Monday ....
    }
  }
  # custom repetition - monthly
  timeslots {
    time_zone  = "Asia/Calcutta"
    start_time = "2022-04-08T06:22:14.975Z"
    end_time   = "2022-04-28T06:22:14.975Z"
    ends_on    = "2022-04-28T06:22:14.975Z"
    repetition = "custom"
    is_allday  = false
    ends_never = true
    custom {
      repeats       = "month"
      repeats_count = 6
    }
  }
}
