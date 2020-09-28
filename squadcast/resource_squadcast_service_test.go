package squadcast

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-provider-squadcast/types"
)

func TestAccSquadcastService_basic(t *testing.T) {
	var test types.Test

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTestCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: interpolateTerraformTemplate(testAccTestConfigBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccTestCheckExists("squadcast_service.pingdom", &test),
				),
			},
		},
	})
}

func testAccTestCheckExists(rn string, test *types.Test) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("TestID not set")
		}

		return nil
	}
}

// TODO: Perform search by service name and if the service exist return error
// rs.Primary.ID will be set to service_name, i.e pingdom_monitoring
func testAccTestCheckDestroy(s *terraform.State) error {
	return nil
}

func interpolateTerraformTemplate(template string) string {
	testContactGroupID := "pingdom_monitoring"

	if v := os.Getenv("SQUADCAST_SERVICE_NAME"); v != "" {
		testContactGroupID = v
	}

	return fmt.Sprintf(template, testContactGroupID)
}

const testAccTestConfigBasic = `
resource "squadcast_service" "pingdom" {
	name = "%s"
	description = "Service created from Terraform acceptance testing"
	escalation_policy_id = "5f35a422ce4a1800086df873"
	email_prefix = "xyz@gmal.com"
}
`
