data "squadcast_team" "example_team" {
  name = "example team name"
}

data "squadcast_user" "example_user" {
  email = "test@example.com"
}

data "squadcast_squad" "example_squad" {
  name = "example squad name"
  team_id = data.squadcast_team.example_team.id
}

data "squadcast_schedule_v2" "example_schedule_v2" {
  name = "example schedule name"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_escalation_policy" "example_escalaion_policy" {
  name        = "example escalation policy name"
  description = "It's an amazing policy"

  team_id = data.squadcast_team.example_team.id

  rules {
    delay_minutes = 0

    targets {
      id   = data.squadcast_user.example_user.id
      type = "user"
    }

    targets {
      id   = data.squadcast_schedule_v2.example_schedule_v2.id
      type = "schedulev2"
    }
  }

  rules {
    delay_minutes = 5

    targets {
      id   = data.squadcast_user.example_user.id
      type = "user"
    }

    targets {
      id   = data.squadcast_squad.example_squad.id
      type = "squad"
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
      id   = data.squadcast_squad.example_squad.id
      type = "squad"
    }

    targets {
      id   = data.squadcast_schedule_v2.example_schedule_v2.id
      type = "schedulev2"
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

  entity_owner {
    id  = data.squadcast_user.example_user.id
    type = "user"
  }
}
