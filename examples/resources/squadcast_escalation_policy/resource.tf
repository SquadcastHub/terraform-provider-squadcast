data "squadcast_team" "example_resource_name" {
  name = "example team name"
}

data "squadcast_user" "example_resource_name" {
  email = "test@example.com"
}

data "squadcast_squad" "example_resource_name" {
  name = "example squad name"
}

data "squadcast_schedule" "example_resource_name" {
  name = "example schedule name"
}

resource "squadcast_escalation_policy" "example_resource_name" {
  name        = "example escalation policy name"
  description = "It's an amazing policy"

  team_id = data.squadcast_team.example_resource_name.id

  rules {
    delay_minutes = 0

    targets {
      id   = data.squadcast_user.example_resource_name.id
      type = "user"
    }

    targets {
      id   = data.squadcast_user.example_resource_name.id
      type = "user"
    }
  }

  rules {
    delay_minutes = 5

    targets {
      id   = data.squadcast_user.example_resource_name.id
      type = "user"
    }

    targets {
      id   = data.squadcast_user.example_resource_name.id
      type = "user"
    }

    notification_channels = ["Phone"]

    repeat {
      times         = 1
      delay_minutes = 5
    }
  }

  rules {
    delay_minutes = 10

    targets {
      id   = data.squadcast_squad.example_resource_name.id
      type = "squad"
    }

    targets {
      id   = data.squadcast_schedule.example_resource_name.id
      type = "schedule"
    }

    round_robin {
      enabled = true

      rotation {
        enabled       = true
        delay_minutes = 1
      }
    }
  }

  repeat {
    times         = 2
    delay_minutes = 10
  }
}