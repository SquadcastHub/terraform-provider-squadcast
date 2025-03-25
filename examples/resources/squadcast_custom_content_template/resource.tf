data "squadcast_team" "example_team" {
  name = "example team name"
}

data "squadcast_service" "example_service" {
  name    = "example service name"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_custom_content_template" "sample2" {
  service_id           = data.squadcast_service.example_service.id
  message_template     = "{{.labels.alertname}}-{{.labels.deployment}}"
  description_template = "{{.labels.description}}"
  alert_source         = "APImetrics"
}
