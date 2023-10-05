data "squadcast_team" "team" {
  name = "Example Team"
}

data "squadcast_service" "service" {
  name = "Example Service"
  team_id = data.squadcast_team.team.id
}

resource "squadcast_ger" "ger" {
  name = "Example GER"
  description = "Example GER Description"
  team_id =   data.squadcast_team.team.id
  entity_owner {
    id = data.squadcast_team.team.id
    type = "team"
  }
}