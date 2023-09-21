data "squadcast_team" "example_team" {
  name = "Example Team"
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
    id = data.squadcast_team.example_team.id
    type = "team"
  }
}

resource "squadcast_ger_ruleset" "example_ger_ruleset" {
    ger_id = squadcast_ger.example_ger.id
    alert_source = "Prometheus"
    catch_all_action = {
        "route_to": data.squadcast_service.example_service.id
    }
}

resource "squadcast_ger_ruleset_rule" "example_ger_ruleset_rule_1" {
    ger_id = squadcast_ger.example_ger.id
    alert_source = squadcast_ger_ruleset.example_ger_ruleset.alert_source
    expression = "alertname == \"DeploymentReplicasNotUpdated\""
    description = "Example GER Ruleset Rule"
    action = {
        "route_to": data.squadcast_service.example_service.id
    }
}

resource "squadcast_ger_ruleset_rule" "example_ger_ruleset_rule_2" {
    ger_id = squadcast_ger.example_ger.id
    alert_source = squadcast_ger_ruleset.example_ger_ruleset.alert_source
    expression = "component == \"kube-state-metrics\""
    description = "Example GER Ruleset Rule"
    action = {
        "route_to": data.squadcast_service.example_service.id
    }
}

resource "squadcast_ger_ruleset_rules_ordering" "rule_ordering" {
    ger_id = squadcast_ger.ger.id
    alert_source = squadcast_ger_ruleset.ger_ruleset_1.alert_source
    ordering = [
        squadcast_ger_ruleset_rule.ger_ruleset_rule_2.id,
        squadcast_ger_ruleset_rule.ger_ruleset_rule_1.id,
    ]
}