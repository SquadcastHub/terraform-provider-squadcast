data "squadcast_team" "example_team" {
  name = "example test name"
}

data "squadcast_service" "example_service" {
  name = "example service name"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_tagging_rules" "example_tagging_rules" {
  team_id    = data.squadcast_team.example_team.id
  service_id = data.squadcast_service.example_service.id

  rules {
    is_basic   = false
    expression = "payload[\"event_id\"] == 40"

    tags {
      key   = "MyTag"
      value = "foo"
      color = "#ababab"
    }
  }

  rules {
    is_basic = true

    basic_expressions {
      lhs = "payload[\"foo\"]"
      op  = "is"
      rhs = "bar"
    }

    tags {
      key   = "MyTag"
      value = "foo"
      color = "#ababab"
    }

    tags {
      key   = "MyTag2"
      value = "bar"
      color = "#f0f0f0"
    }
  }
}

# addTags must be set in expression when tags are not passed
resource "squadcast_tagging_rules" "example_tagging_rules_resource_withouttags" {
  team_id    = data.squadcast_team.example_team.id
  service_id = data.squadcast_service.example_service.id

  rules {
    is_basic   = false
    expression = "addTag(\"EventType\", payload.details.event_type_key, \"#037916\")"
  }
}
