data "squadcast_user" "example_user" {
  email = "test@example.com"
}

data "squadcast_team" "example_team" {
  name = "example team name"
}

data "squadcast_service" "example_service" {
  name = "example service name"
  team_id = data.squadcast_team.example_team.id
}

data "squadcast_service" "example_service_new" {
  name = "example service name 2"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_workflow" "example_workflow_with_simplest_filter" {
   title = "test workflow"
   description = "Test workflow description"
   owner_id = data.squadcast_team.example_team.id
   enabled = true
   trigger = "incident_triggered"
   filters {
      filters {
        type = "priority_is"
        value = "P1"
      }
   }
   entity_owner {
      type = "user" 
      id = data.squadcast_user.example_user.id
   }
   tags {
      key = "tagKey"
      value = "tagValue"
      color = "#000000"
   }
}

resource "squadcast_workflow" "example_workflow_with_simple_filters" {
   title = "test workflow"
   description = "Test workflow description"
   owner_id = data.squadcast_team.example_team.id
   enabled = true
   trigger = "incident_triggered"
   filters {
      condition = "or"
      filters {
        type = "priority_is"
        value = "P1"
      }
      filters {
         type = "priority_is"
         value = "UNSET"
      }
   }
   entity_owner {
      type = "user" 
      id = data.squadcast_user.example_user.id
   }
   tags {
      key = "tagKey"
      value = "tagValue"
      color = "#000000"
   }
}

resource "squadcast_workflow" "example_workflow_with_advanced_filters" {
   title = "test workflow"
   description = "Test workflow description"
   owner_id = data.squadcast_team.example_team.id
   enabled = true
   trigger = "incident_triggered"
   filters {
      condition = "or"
      filters {
         condition = "and"
         filters {
            type = "tag_is"
            key = "hello"
            value = "world"            
         }         
         filters {
            type = "tag_is"
            key = "service"
            value = "payment-gw"            
         }
      }
      filters {
        condition = "or"
        filters {
          type = "service_is"
          value = data.squadcast_service.example_service.id
        }
        filters {
          type = "service_is"
          value = data.squadcast_service.example_service_new.id
        }
      }
      filters {
         type = "priority_is"
         value = "P1"
      }
   }
   entity_owner {
      type = "user" 
      id = data.squadcast_user.example_user.id
   }
   tags {
      key = "tagKey"
      value = "tagValue"
      color = "#000000"
   }
}