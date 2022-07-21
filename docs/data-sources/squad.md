---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "squadcast_squad Data Source - terraform-provider-squadcast"
subcategory: ""
description: |-
  Squads https://support.squadcast.com/docs/squads are smaller groups of members within Teams. Squads could correspond to groups of people that are responsible for specific projects within a Team.Use this data source to get information about a specific Squad.
---

# squadcast_squad (Data Source)

[Squads](https://support.squadcast.com/docs/squads) are smaller groups of members within Teams. Squads could correspond to groups of people that are responsible for specific projects within a Team.Use this data source to get information about a specific Squad.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the Squad.
- `team_id` (String) Team id.

### Read-Only

- `id` (String) Squad id.
- `member_ids` (List of String)

