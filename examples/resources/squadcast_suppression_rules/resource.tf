data "squadcast_team" "example_team" {
  name = "exammple team name"
}

data "squadcast_service" "example_service" {
  name    = "example service name"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_suppression_rules" "example_suppression_rules" {
  team_id    = data.squadcast_team.example_team.id
  service_id = data.squadcast_service.example_service.id

  rules {
    is_basic    = false
    description = "not basic"
    expression  = "payload[\"event_id\"] == 40"
  }
}

resource "squadcast_suppression_rules" "example_time_based_suppression_rules" {
  team_id    = data.squadcast_team.example_team.id
  service_id = data.squadcast_service.example_service.id

  rules {
    is_basic     = false
    description  = "not basic"
    expression   = "payload[\"event_id\"] == 40"
    is_timebased = true
    timeslots {
      time_zone  = "Asia/Calcutta"
      start_time = "2022-04-08T06:22:14.975Z"
      end_time   = "2022-04-28T06:22:14.975Z"
      ends_on    = "2022-04-28T06:22:14.975Z"
      repetition = "custom"
      is_allday  = false
      ends_never = true
      is_custom  = true
      custom {
        repeats             = "day"
        repeats_count       = 2
        repeats_on_weekdays = [0, 1] # 0 - Sunday, 1 - Monday ...
      }
    }
  }
}
