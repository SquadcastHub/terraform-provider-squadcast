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

resource "squadcast_workflow_action" "example_workflow1" {
   workflow_id = squadcast_workflow.testing_workflows.id
   name = "slack_create_incident_channel" 
   auto_name = false
   channel_name = "enter-channel-name"
}

resource "squadcast_workflow_action" "example_workflow2" {
   workflow_id = squadcast_workflow.testing_workflows.id
   name = "slack_archive_channel" 
}

resource "squadcast_workflow_action" "example_workflow3" {
   workflow_id = squadcast_workflow.testing_workflows.id
   name = "sq_update_incident_priority"
   priority = "P2"
}

resource "squadcast_workflow_action" "example_workflow4" {
   workflow_id = squadcast_workflow.testing_workflows.id
   name = "sq_add_communication_channel"
   channels {
      type = "chat_room"
      link = "https://chat.squadcast.com/room/123456"
      display_text = "enter your display text here"
   }
}

resource "squadcast_workflow_action" "example_workflow5" {
   workflow_id = squadcast_workflow.testing_workflows.id
   name = "sq_mark_incident_slo_affecting"
   slo = "2119"
   slis = ["errors"]
}

resource "squadcast_workflow_action_ordering" "def"{
   workflow_id = squadcast_workflow.example_workflow.id
   action_order = [squadcast_workflow_action.example_workflow5.id, squadcast_workflow_action.example_workflow3.id, 
        squadcast_workflow_action.example_workflow4.id, squadcast_workflow_action.example_workflow1.id, squadcast_workflow_action.example_workflow2.id]
}