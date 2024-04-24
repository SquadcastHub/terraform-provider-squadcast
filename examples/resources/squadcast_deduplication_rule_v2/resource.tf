data "squadcast_team" "example_team" {
  name = "example team name"
}

data "squadcast_service" "example_service" {
  name = "example service name"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_deduplication_rule_v2" "example_deduplication_rule" {
    service_id = data.squadcast_service.example_service.id
    is_basic    = false
    description = "not basic"
    expression  = "payload[\"event_id\"] == 40"
}

resource "squadcast_deduplication_rule_v2" "example_basic_deduplication_rule" {
    service_id = data.squadcast_service.example_service.id
    is_basic    = true
    description = "basic"
    dependency_deduplication = true
    time_window = 2
    time_unit = "hour"

    basic_expressions {
        lhs = "payload[\"foo\"]"
        op  = "is"
        rhs = "bar"
    }
}
