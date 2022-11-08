data "squadcast_team" "default" {
  name = "example test name"
}

data "squadcast_service" "default" {
  name = "example service name"
}

resource "squadcast_tagging_rules" "default" {
  team_id    = data.squadcast_team.default.id
  service_id = data.squadcast_service.default.id

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
resource "squadcast_tagging_rules" "withouttags" {
  team_id    = data.squadcast_team.default.id
  service_id = data.squadcast_service.default.id

  rules {
    is_basic   = false
    expression = "addTag(\"EventType\", payload.details.event_type_key, \"#037916\")"
  }
}