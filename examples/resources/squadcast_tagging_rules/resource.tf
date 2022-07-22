data "squadcast_team" "example" {
  name = "test"
}

data "squadcast_service" "example" {
  name = "test-parent"
}

resource "squadcast_tagging_rules" "test" {
  team_id    = data.squadcast_team.example.id
  service_id = data.squadcast_service.example.id

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