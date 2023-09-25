---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "squadcast_ger Resource - terraform-provider-squadcast"
subcategory: ""
description: |-
  Global Event Ruleset (GER) is a centralized set of rules that defines service routes for incoming events.
---

# squadcast_ger (Resource)

Global Event Ruleset (GER) is a centralized set of rules that defines service routes for incoming events.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `entity_owner` (Block List, Min: 1, Max: 1) GER owner. (see [below for nested schema](#nestedblock--entity_owner))
- `name` (String) GER name.
- `team_id` (String) Team id.

### Optional

- `description` (String) GER description.

### Read-Only

- `id` (String) GER id.
- `routing_key` (String) Routing Key is an identifier used to determine the ruleset that an incoming event belongs to. It is a common key that associates multiple alert sources with their configured rules, ensuring events are routed to the appropriate services when the defined criteria are met.

<a id="nestedblock--entity_owner"></a>
### Nested Schema for `entity_owner`

Required:

- `id` (String) GER owner id.
- `type` (String) GER owner type. (user or squad or team)

## Import

Import is supported using the following syntax:

```shell
# gerID

terraform import squadcast_ger.example_ger_import "53"
```