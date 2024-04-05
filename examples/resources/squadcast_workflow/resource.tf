resource "squadcast_workflow" "example_workflow" {
   title = "test workflow"
   description = "Test workflow description"
   owner_id = "63bfabae865e9c93cd31756e"
   enabled = true
   trigger = "incident_triggered"
   # TODO: Needs to accomodate to the new structure
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