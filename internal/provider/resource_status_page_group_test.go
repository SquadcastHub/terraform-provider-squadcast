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

func TestAccResourceStatusPageGroup(t *testing.T) {
	statusPageGroupName := acctest.RandomWithPrefix("statusPageGroup")

	resourceName := "squadcast_status_page_group.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckStatusPageGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStatusPageGroupConfig(statusPageGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "status_page_id", "100"),
					resource.TestCheckResourceAttr(resourceName, "name", statusPageGroupName),
				),
			},
			{
				Config: testAccResourceStatusPageGroupConfig_update(statusPageGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "status_page_id", "100"),
					resource.TestCheckResourceAttr(resourceName, "name", statusPageGroupName),
				),
			},
		},
	})
}

func testAccCheckStatusPageGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_status_page_group" {
			continue
		}

		_, err := client.GetStatusPageGroupById(context.Background(), rs.Primary.Attributes["status_page_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("expected status page to be destroyed, %s found", rs.Primary.ID)
		}

		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceStatusPageGroupConfig(statusPageGroupName string) string {
	return fmt.Sprintf(`
resource "squadcast_status_page_group" "test" {
	status_page_id = "100"
	name = "%s"
}
	`, statusPageGroupName)
}

func testAccResourceStatusPageGroupConfig_update(statusPageGroupName string) string {
	return fmt.Sprintf(`
resource "squadcast_status_page_group" "test" {
	status_page_id = "100"
	name = "%s"
}
	`, statusPageGroupName)
}
