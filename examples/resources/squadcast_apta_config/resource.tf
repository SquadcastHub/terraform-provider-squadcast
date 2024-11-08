data "squadcast_team" "example_team" {
  name = "example team name"
}

data "squadcast_service" "example_service" {
  name    = "example service name"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_apta_config" "apta" {
  service_id = data.squadcast_service.example_service.id
  is_enabled = true
  timeout    = 5 # 2, 3, 5, 10, 15
}
