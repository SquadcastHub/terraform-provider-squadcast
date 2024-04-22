resource "squadcast_workflow" "example_workflow" {
   title = "test workflow"
   description = "Test workflow description"
   owner_id = "63bfabae865e9c93cd31756e"
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
      id = "63209531af0f36245bfac82f"
   }
   tags {
      key = "tagKey"
      value = "tagValue"
      color = "#000000"
   }
}


resource "squadcast_workflow" "example_workflow_with_advaced_filters" {
   title = "test workflow"
   description = "Test workflow description"
   owner_id = "63bfabae865e9c93cd31756e"
   enabled = true
   trigger = "incident_triggered"
   filters {
      condition = "or"
      filters {
         condition = "and"
         filters {
            type = "tag_is"
            key = "hello1"
            value = "world1"            
         }         
         filters {
            type = "tag_is"
            key = "hello"
            value = "world"            
         }
      }
      filters {
         type = "priority_is"
         value = "P3"
      }
   }
   entity_owner {
      type = "user" 
      id = "63209531af0f36245bfac82f"
   }
   tags {
      key = "tagKey"
      value = "tagValue"
      color = "#000000"
   }
}