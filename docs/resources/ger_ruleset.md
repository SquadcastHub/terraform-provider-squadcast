---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "squadcast_ger_ruleset Resource - terraform-provider-squadcast"
subcategory: ""
description: |-
  GER Ruleset is a set of rules and configurations in Squadcast. It allows users to define how alerts are routed to services without the need to set up individual webhooks for each alert source.
---

# squadcast_ger_ruleset (Resource)

GER Ruleset is a set of rules and configurations in Squadcast. It allows users to define how alerts are routed to services without the need to set up individual webhooks for each alert source.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `alert_source` (String) An alert source refers to the origin of an event (alert), such as a monitoring tool. These alert sources are associated with specific rules in GER, determining where events from each source should be routed. Find all alert sources supported on Squadcast [here](https://www.squadcast.com/integrations).
- `ger_id` (String) GER id.

### Optional

- `catch_all_action` (Map of String) The "Catch-All Action", when configured, specifies a fall back service. If none of the defined rules for an incoming event evaluate to true, the incoming event is routed to the Catch-All service, ensuring no events are missed.

### Read-Only

- `alert_source_shortname` (String) Shortname of the linked alert source.
- `alert_source_version` (String) Version of the linked alert source.
- `id` (String) GER Ruleset id.

## Import

Import is supported using the following syntax:

```shell
# gerID:alertSourceName
terraform import squadcast_ger_ruleset.example_ger_ruleset_import "53:Grafana"
```