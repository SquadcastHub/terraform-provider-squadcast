package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
)

func TestAccResourceAPTAConfig(t *testing.T) {
	resourceName := "squadcast_apta_config.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAPTAConfigDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic configuration
				Config: testAccResourceAPTAConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "service_id", "61361611c2fc70c3101ca7dd"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "5"),
				),
			},
			{
				// Test import
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "61361611c2fc70c3101ca7dd",
			},
			{
				// Test update with different timeout
				Config: testAccResourceAPTAConfigDifferentTimeout(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "service_id", "61361611c2fc70c3101ca7dd"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
				),
			},
			{
				// Test disabled configuration
				Config: testAccResourceAPTAConfigDisabled(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "service_id", "61361611c2fc70c3101ca7dd"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "5"),
				),
			},
		},
	})
}

func testAccCheckAPTAConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_apta_config" {
			continue
		}

		_, err := client.GetServiceById(context.Background(), "", rs.Primary.Attributes["service_id"])
		if err != nil {
			if api.IsResourceNotFoundError(err) {
				return nil
			}
			return err
		}

		return fmt.Errorf("APTA Config still exists")
	}

	return nil
}

func testAccResourceAPTAConfigBasic() string {
	return fmt.Sprintf(`
resource "squadcast_apta_config" "test" {
	service_id = "61361611c2fc70c3101ca7dd"
	is_enabled = true
	timeout    = 5
}
`)
}

func testAccResourceAPTAConfigDifferentTimeout() string {
	return fmt.Sprintf(`
resource "squadcast_apta_config" "test" {
	service_id = "61361611c2fc70c3101ca7dd"
	is_enabled = true
	timeout    = 10
}
`)
}

func testAccResourceAPTAConfigDisabled() string {
	return fmt.Sprintf(`
resource "squadcast_apta_config" "test" {
	service_id = "61361611c2fc70c3101ca7dd"
	is_enabled = false
	timeout    = 5
}
`)
}
