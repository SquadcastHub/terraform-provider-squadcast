package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
)

func TestAccResourceIAGConfig(t *testing.T) {
	resourceName := "squadcast_iag_config.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckIAGConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIAGConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "service_id", "61361611c2fc70c3101ca7dd"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "grouping_window", "5"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "61361611c2fc70c3101ca7dd", // service_id
			},
		},
	})
}

func testAccCheckIAGConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_iag_config" {
			continue
		}

		_, err := client.GetServiceById(context.Background(), "", rs.Primary.Attributes["service_id"])
		if err != nil {
			if api.IsResourceNotFoundError(err) {
				return nil
			}
			return err
		}

		return fmt.Errorf("IAG Config still exists")
	}

	return nil
}

func testAccResourceIAGConfigBasic() string {
	return fmt.Sprintf(`
resource "squadcast_iag_config" "test" {
	service_id      = "61361611c2fc70c3101ca7dd"
	is_enabled      = true
	grouping_window = 5
}
`)
}
