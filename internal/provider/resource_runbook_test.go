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

func TestAccResourceRunbook(t *testing.T) {
	runbookName := acctest.RandomWithPrefix("runbook")

	resourceName := "squadcast_runbook.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckRunbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRunbookConfig(runbookName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", runbookName),
					resource.TestCheckResourceAttr(resourceName, "steps.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "steps.0.content", "some text here"),
				),
			},
			{
				Config: testAccResourceRunbookConfig_update(runbookName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", runbookName),
					resource.TestCheckResourceAttr(resourceName, "steps.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "steps.0.content", "some text here"),
					resource.TestCheckResourceAttr(resourceName, "steps.1.content", "some text here 2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "613611c1eb22db455cfa789f:" + runbookName,
			},
		},
	})
}

func testAccCheckRunbookDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_runbook" {
			continue
		}

		_, err := client.GetRunbookById(context.Background(), rs.Primary.Attributes["team_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("expected runbook to be destroyed, %s found", rs.Primary.ID)
		}

		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceRunbookConfig(runbookName string) string {
	return fmt.Sprintf(`
resource "squadcast_runbook" "test" {
	name = "%s"
	team_id = "613611c1eb22db455cfa789f"

	steps {
		content = "some text here"
	}
}
	`, runbookName)
}

func testAccResourceRunbookConfig_update(runbookName string) string {
	return fmt.Sprintf(`
resource "squadcast_runbook" "test" {
	name = "%s"
	team_id = "613611c1eb22db455cfa789f"

	steps {
		content = "some text here"
	}

	steps {
		content = "some text here 2"
	}
}
	`, runbookName)
}
