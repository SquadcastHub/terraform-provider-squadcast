data "squadcast_team" "example_team" {
  name = "Example Team"
}

data "squadcast_user" "user" {
  email = "john@example.com"
}


data "squadcast_service" "example_service" {
  name = "Example Service"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_ger" "example_ger" {
  name = "Example GER"
  description = "Example GER Description"
  team_id =   data.squadcast_team.example_team.id
  entity_owner {
    id = data.squadcast_user.user.id
    type = "user"
  }
}

resource "squadcast_ger_ruleset" "example_ger_ruleset" {
    ger_id = squadcast_ger.example_ger.id
    alert_source = "Prometheus"
    catch_all_action = {
        "route_to": data.squadcast_service.example_service.id
    }
}