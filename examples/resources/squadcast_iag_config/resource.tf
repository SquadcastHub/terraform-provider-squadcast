data "squadcast_team" "example_team" {
  name = "example team name"
}

data "squadcast_service" "example_service" {
  name    = "example service name"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_iag_config" "apta" {
  service_id      = data.squadcast_service.example_service.id
  is_enabled      = true
  grouping_window = 5 # 5, 10, 15, 20, 45, 60, 60, 240, 480, 720, 1440
}
