data "squadcast_team" "example_team" {
  name = "example team name"
}

data "squadcast_service" "example_service" {
  name = "example service name"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_service_maintenance" "example_service_maintenance" {
  service_id = data.squadcast_service.example_service.id

  windows {
    from             = "2032-06-01T10:30:00.000Z"
    till             = "2032-06-01T11:30:00.000Z"
    repeat_till      = "2032-06-30T10:30:00.000Z"
    repeat_frequency = "week"
  }

  windows {
    from = "2032-07-01T10:30:00.000Z"
    till = "2032-07-02T10:30:00.000Z"
  }
}
