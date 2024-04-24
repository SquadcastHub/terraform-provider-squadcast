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

resource "squadcast_workflow_action" "example_workflow" {
   workflow_id = squadcast_workflow.testing_workflows.id
   name = "slack_create_incident_channel" 
   auto_name = false
   channel_name = "enter-channel-name"
}

resource "squadcast_workflow_action" "example_workflow" {
   workflow_id = squadcast_workflow.testing_workflows.id
   name = "slack_archive_channel" 
}

resource "squadcast_workflow_action" "example_workflow" {
   workflow_id = squadcast_workflow.testing_workflows.id
   name = "slack_message_channel"
   channel_id = "C06P4473BJA"
   message = "test incident created..."
}

resource "squadcast_workflow_action" "example_workflow" {
   workflow_id = squadcast_workflow.testing_workflows.id
   name = "sq_trigger_manual_webhook" 
   webhook_id = "660edb863a1cefa8f291aebe"
}

resource "squadcast_workflow_action" "example_workflow" {
   workflow_id = squadcast_workflow.testing_workflows.id
   name = "sq_send_email"
   to = ["abc@squadcast.com", "xyz@squadcast.com"]
   subject = "enter your subject here"
   body = "enter your body here"
}

resource "squadcast_workflow_action" "example_workflow" {
   workflow_id = squadcast_workflow.testing_workflows.id
   name = "sq_make_http_call"
   url = "https://httpbin.org/post"
   method = "GET"
   headers {
      key = "content-type"
      value = "application/json"
   }
   body = "{\"key\":\"value\"}"
}

resource "squadcast_workflow_action" "example_workflow" {
   workflow_id = squadcast_workflow.testing_workflows.id
   name = "sq_update_incident_priority"
   priority = "P2"
}

resource "squadcast_workflow_action" "example_workflow" {
   workflow_id = squadcast_workflow.testing_workflows.id
   name = "sq_add_communication_channel"
   channels {
      type = "chat_room"
      link = "https://chat.squadcast.com/room/123456"
      display_text = "enter your display text here"
   }
}

resource "squadcast_workflow_action" "example_workflow" {
   workflow_id = squadcast_workflow.testing_workflows.id
   name = "sq_mark_incident_slo_affecting"
   slo = "2119"
   slis = ["errors"]
}

resource "squadcast_workflow_action" "example_workflow" {
   workflow_id = squadcast_workflow.testing_workflows.id
   name = "sq_attach_runbooks"
   runbooks = ["660ced558d1d4df4a61823ee", "660d46f62f8acc7786618202"]
}