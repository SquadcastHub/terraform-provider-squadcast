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

func TestAccResourceStatusPageComponent(t *testing.T) {
	statusPageComponentName := acctest.RandomWithPrefix("statusPageComponent")

	resourceName := "squadcast_status_page_component.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckStatusPageComponentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStatusPageComponentConfig(statusPageComponentName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "status_page_id", "100"),
					resource.TestCheckResourceAttr(resourceName, "name", statusPageComponentName),
					resource.TestCheckResourceAttr(resourceName, "description", "Sample status page component description."),
				),
			},
			{
				Config: testAccResourceStatusPageComponentConfig_update(statusPageComponentName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "status_page_id", "100"),
					resource.TestCheckResourceAttr(resourceName, "name", statusPageComponentName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated Sample status page component description."),
				),
			},
			{
				Config: testAccResourceStatusPageComponentConfig_group(statusPageComponentName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "status_page_id", "100"),
					resource.TestCheckResourceAttr(resourceName, "name", statusPageComponentName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated Sample status page component description."),
					resource.TestCheckResourceAttr(resourceName, "group_id", "200"),
				),
			},
		},
	})
}

func testAccCheckStatusPageComponentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_status_page_component" {
			continue
		}

		_, err := client.GetStatusPageComponentById(context.Background(), rs.Primary.Attributes["status_page_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("expected status page to be destroyed, %s found", rs.Primary.ID)
		}

		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceStatusPageComponentConfig(statusPageComponentName string) string {
	return fmt.Sprintf(`
resource "squadcast_status_page_component" "test" {
	status_page_id = "100"
	name = "%s"
	description = "Sample status page component description."
}
	`, statusPageComponentName)
}

func testAccResourceStatusPageComponentConfig_update(statusPageComponentName string) string {
	return fmt.Sprintf(`
resource "squadcast_status_page_component" "test" {
	status_page_id = "100"
	name = "%s"
	description = "Updated Sample status page component description."
}
	`, statusPageComponentName)
}

func testAccResourceStatusPageComponentConfig_group(statusPageComponentName string) string {
	return fmt.Sprintf(`
resource "squadcast_status_page_component" "test" {
	status_page_id = "100"
	name = "%s"
	description = "Updated Sample status page component description."
	group_id = "200"
}
	`, statusPageComponentName)
}
