data "squadcast_team" "example_team_resource" {
  name = "example team name"
}

data "squadcast_user" "example_user_resource" {
  email = "test@example.com"
}

data "squadcast_service" "example_service_resource" {
  name = "example service name"
  team_id = data.squadcast_team.example_team_resource.id
}

resource "squadcast_slo" "example_slo_resource" {
  name               = "checkout-flow"
  description        = "Slo for checkout flow"
  target_slo         = 99.99
  service_ids        = [data.squadcast_service.example_service_resource.id]
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
    user_ids = [data.squadcast_user.example_user_resource.id]
  }

  team_id = data.squadcast_team.example_team_resource.id
}