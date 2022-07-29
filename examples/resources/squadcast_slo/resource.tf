data "squadcast_team" "example_resource_name" {
  name = "example team name"
}

data "squadcast_user" "example_resource_name" {
  email = "test@example.com"
}

data "squadcast_service" "example_resource_name" {
  name = "example service name"
}

resource "squadcast_slo" "example_resource_name" {
  name               = "checkout-flow"
  description        = "Slo for checkout flow"
  target_slo         = 99.99
  service_ids        = [data.squadcast_service.example_resource_name.id]
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
    user_ids = [data.squadcast_user.example_resource_name.id]
  }

  team_id = data.squadcast_team.example_resource_name.id
}