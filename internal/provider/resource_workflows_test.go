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

func TestAccResourceWorkflows(t *testing.T) {
	workflowTitle := acctest.RandomWithPrefix("test-workflow")
	resourceName := "squadcast_workflows.test_workflows"
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckWorkflowsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowsConfig(workflowTitle),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "title", workflowTitle),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
				),
			},
			{
				Config: testAccResourceWorkflows_update(workflowTitle),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "title", workflowTitle),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
				),
			},
		},
	})
}

func testAccCheckWorkflowsDestroy(s *terraform.State) error {
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

func testAccResourceWorkflowsConfig(workflowTitle string) string {
	return fmt.Sprintf(`
	resource "squadcast_workflows" "test_workflows" {
		title = "%s"
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
`, workflowTitle)
}

func testAccResourceWorkflows_update(workflowTitle string) string {
	return fmt.Sprintf(`
	resource "squadcast_workflows" "test_workflows" {
		title = "%s"
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
`, workflowTitle)
}
