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

func TestAccResourceStatusPage(t *testing.T) {
	statusPageName := acctest.RandomWithPrefix("statusPage")

	resourceName := "squadcast_status_page.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckStatusPageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStatusPageConfig(statusPageName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", statusPageName),
					resource.TestCheckResourceAttr(resourceName, "description", "Sample status page description."),
					resource.TestCheckResourceAttr(resourceName, "is_public", "true"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "sq-statuspage"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Asia/Kolkata"),
					resource.TestCheckResourceAttr(resourceName, "contact_email", "test@squadcast.com"),
					resource.TestCheckResourceAttr(resourceName, "owner.0.id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "owner.0.type", "team"),
					resource.TestCheckResourceAttr(resourceName, "theme_color.0.primary", "#000000"),
					resource.TestCheckResourceAttr(resourceName, "theme_color.0.secondary", "#ffffff"),
					resource.TestCheckResourceAttr(resourceName, "allow_webhook_subscription", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_components_subscription", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_maintenance_subscription", "true"),
				),
			},
			{
				Config: testAccResourceStatusPageConfig_update(statusPageName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", statusPageName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated Sample status page description."),
					resource.TestCheckResourceAttr(resourceName, "is_public", "false"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "sq-statuspage"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Asia/Kolkata"),
					resource.TestCheckResourceAttr(resourceName, "contact_email", "contact@squadcast.com"),
					resource.TestCheckResourceAttr(resourceName, "owner.0.id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "owner.0.type", "team"),
					resource.TestCheckResourceAttr(resourceName, "theme_color.0.primary", "#000000"),
					resource.TestCheckResourceAttr(resourceName, "theme_color.0.secondary", "#ffffff"),
					resource.TestCheckResourceAttr(resourceName, "allow_webhook_subscription", "false"),
					resource.TestCheckResourceAttr(resourceName, "allow_components_subscription", "false"),
					resource.TestCheckResourceAttr(resourceName, "allow_maintenance_subscription", "false"),
				),
			},
		},
	})
}

func testAccCheckStatusPageDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_status_page" {
			continue
		}

		_, err := client.GetStatusPageById(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("expected status page to be destroyed, %s found", rs.Primary.ID)
		}

		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceStatusPageConfig(statusPageName string) string {
	return fmt.Sprintf(`
resource "squadcast_status_page" "test" {
	name = "%s"
	team_id = "613611c1eb22db455cfa789f"
	description = "Sample status page description."
	is_public = true
	domain_name = "sq-statuspage"
	timezone = "Asia/Kolkata"
	contact_email = "test@squadcast.com"
	owner {
		id = "613611c1eb22db455cfa789f"
		type = "team"
	}
	theme_color {
		primary = "#000000"
		secondary = "#ffffff"
	}
	allow_webhook_subscription = true
	allow_components_subscription = true
	allow_maintenance_subscription = true
}
	`, statusPageName)
}

func testAccResourceStatusPageConfig_update(statusPageName string) string {
	return fmt.Sprintf(`
resource "squadcast_status_page" "test" {
	name = "%s"
	team_id = "613611c1eb22db455cfa789f"
	description = "Updated Sample status page description."
	is_public = false
	domain_name = "sq-statuspage"
	timezone = "Asia/Kolkata"
	contact_email = "contact@squadcast.com"
	owner {
		id = "613611c1eb22db455cfa789f"
		type = "team"
	}
	theme_color {
		primary = "#000000"
		secondary = "#ffffff"
	}
	allow_webhook_subscription = false
	allow_components_subscription = false
	allow_maintenance_subscription = false
}
	`, statusPageName)
}
