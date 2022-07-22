data "squadcast_team" "example" {
  name = "test"
}

data "squadcast_user" "example" {
  email = "test@example.com"
}

data "squadcast_service" "example" {
  name = "test-parent"
}

resource "squadcast_slo" "test" {
  name               = "checkout-flow"
  description        = "Slo for checkout flow"
  target_slo         = 99.99
  service_ids        = [data.squadcast_service.example.id]
  slis               = ["latency", "high-err-rate"]
  time_interval_type = "rolling"
  duration_in_days   = 7

  rules {
    name = "breached_error_budget"
  }

  rules {
    name      = "remaining_error_budget"
    threshold = 11
  }

  rules {
    name      = "unhealthy_slo"
    threshold = 1
  }

  notify {
    users = [data.squadcast_user.example.id]
  }

  owner_type = "team"
  team_id    = data.squadcast_team.example.id
}