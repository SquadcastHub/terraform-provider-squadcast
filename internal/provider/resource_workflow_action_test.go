package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
)

func TestAccResourceWorkflowAction(t *testing.T) {
	workflowTitle := acctest.RandomWithPrefix("test-workflow-action")
	resourceName := "squadcast_workflow_action.test_workflow_action"
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckWorkflowActionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowActionConfig(workflowTitle),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "note", "testing workflow action"),
				),
			},
			{
				Config: testAccResourceWorkflowAction_update(workflowTitle),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "note", "testing update workflow action"),
				),
			},
		},
	})
}

func testAccCheckWorkflowActionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_workflows" {
			continue
		}

		_, err := client.GetWorkflowById(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("workflow still exists")
		}
	}

	return nil
}

func testAccResourceWorkflowActionConfig(workflowTitle string) string {
	return fmt.Sprintf(`

	resource "squadcast_workflow_action" "test_workflow_action" {
		workflow_id = squadcast_workflows.test_workflows.id
		name = "sq_add_incident_note"
		note = "testing workflow action"
	}

	resource "squadcast_workflows" "test_workflows" {
		title = "%s"
		description = "Test workflow description"
		owner_id = "63bfabae865e9c93cd31756e"
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
			id = "63209531af0f36245bfac82f"
		}
		tags {
			key = "tagKey"
			value = "tagValue"
			color = "#000000"
		}
	}
`, workflowTitle)
}

func testAccResourceWorkflowAction_update(workflowTitle string) string {
	return fmt.Sprintf(`

	resource "squadcast_workflow_action" "test_workflow_action" {
		workflow_id = squadcast_workflows.test_workflows.id
		name = "sq_add_incident_note"
		note = "testing update workflow action"
	}

	resource "squadcast_workflows" "test_workflows" {
		title = "%s"
		description = "Test workflow description"
		owner_id = "63bfabae865e9c93cd31756e"
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
			id = "63209531af0f36245bfac82f"
		}
		tags {
			key = "tagKey"
			value = "tagValue"
			color = "#000000"
		}
	}
`, workflowTitle)
}
