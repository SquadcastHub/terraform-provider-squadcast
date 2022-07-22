data "squadcast_team" "example" {
  name = "test"
}

data "squadcast_user" "example" {
  email = "test@example.com"
}

data "squadcast_squad" "example" {
  name = "test"
}

data "squadcast_schedule" "example" {
  name = "test"
}

resource "squadcast_escalation_policy" "test" {
  name        = "test escalation policy"
  description = "It's an amazing policy"

  team_id = data.squadcast_team.example.id

  rules {
    delay_minutes = 0

    targets {
      id   = data.squadcast_user.example.id
      type = "user"
    }

    targets {
      id   = data.squadcast_user.example.id
      type = "user"
    }
  }

  rules {
    delay_minutes = 5

    targets {
      id   = data.squadcast_user.example.id
      type = "user"
    }

    targets {
      id   = data.squadcast_user.example.id
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
      id   = data.squadcast_squad.example.id
      type = "squad"
    }

    targets {
      id   = data.squadcast_schedule.example.id
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