resource "squadcast_slo" "test" {
  name               = "checkout-flow"
  description        = "Slo for checkout flow"
  target_slo         = 99.99
  service_ids        = ["615d3e23aff6885f46d291be"]
  slis               = ["latency", "high-err-rate"]
  time_interval_type = "rolling"
  duration_in_days   = 7
  org_id             = "604592dabc35ea0008bb0584"

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
    users = ["5e1c2309342445001180f9c2", "617793e650d38001057faaaf"]
  }

  owner_type = "team"
  team_id    = "611262fcd5b4ea846b534a8a"
}