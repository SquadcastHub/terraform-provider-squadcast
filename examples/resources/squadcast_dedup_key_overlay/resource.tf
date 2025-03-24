data "squadcast_team" "example_team" {
  name = "example team name"
}

data "squadcast_service" "example_service" {
  name    = "example service name"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_dedup_key_overlay" "example" {
  service_id                 = data.squadcast_service.example_service.id
  dedup_key_overlay_template = "Alertname: (?P<alertname>.*)|Summary: (?P<summary>.*)$"
  duration                   = 100
  alert_source               = "APImetrics"
}
