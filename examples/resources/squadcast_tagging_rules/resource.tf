data "squadcast_team" "example_resource_name" {
  name = "example test name"
}

data "squadcast_service" "example_resource_name" {
  name = "example service name"
}

resource "squadcast_tagging_rules" "example_resource_name" {
  team_id    = data.squadcast_team.example_resource_name.id
  service_id = data.squadcast_service.example_resource_name.id

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