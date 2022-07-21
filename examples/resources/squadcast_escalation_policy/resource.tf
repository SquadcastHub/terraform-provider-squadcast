resource "squadcast_escalation_policy" "test" {
  name        = "test escalation policy"
  description = "It's an amazing policy"

  team_id = "owner_id"

  rules {
    delay_minutes = 0

    targets {
      id   = "user_id"
      type = "user"
    }

    targets {
      id   = "user_id"
      type = "user"
    }
  }

  rules {
    delay_minutes = 5

    targets {
      id   = "user_id"
      type = "user"
    }

    targets {
      id   = "user_id"
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
      id   = "squad_id"
      type = "squad"
    }

    targets {
      id   = "schedule_id"
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