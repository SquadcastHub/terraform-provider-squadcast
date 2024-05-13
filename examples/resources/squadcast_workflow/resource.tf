data "squadcast_user" "example_user" {
  email = "test@example.com"
}

data "squadcast_team" "example_team" {
  name = "example team name"
}

resource "squadcast_workflow" "example_workflow_with_simple_filters" {
   title = "test workflow"
   description = "Test workflow description"
   owner_id = data.squadcast_team.example_team.id
   enabled = true
   trigger = "incident_triggered"
   filters {
      fields {
         value = "P1"
      }
      type = "priority_is"
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