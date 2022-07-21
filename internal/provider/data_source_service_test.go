package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceService(t *testing.T) {
	serviceName := acctest.RandomWithPrefix("service")

	resourceName := "data.squadcast_service.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDataSourceConfig(serviceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", serviceName),
					resource.TestCheckResourceAttr(resourceName, "description", "some description here."),
					resource.TestCheckResourceAttr(resourceName, "escalation_policy_id", "61361415c2fc70c3101ca7db"),
					resource.TestCheckResourceAttr(resourceName, "email_prefix", serviceName),
					resource.TestCheckResourceAttrSet(resourceName, "api_key"),
					resource.TestCheckResourceAttr(resourceName, "email", serviceName+"@squadcast.incidents.squadcast.com"),
					resource.TestCheckResourceAttr(resourceName, "dependencies.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "dependencies.0", "squadcast_service.test_parent", "id"),
					resource.TestCheckResourceAttr(resourceName, "alert_source_endpoints.email", serviceName+"@squadcast.incidents.squadcast.com"),
				),
			},
		},
	})
}

func testAccServiceDataSourceConfig(serviceName string) string {
	return fmt.Sprintf(`
resource "squadcast_service" "test_parent" {
	name = "%s-parent"
	team_id = "613611c1eb22db455cfa789f"
	escalation_policy_id = "61361415c2fc70c3101ca7db"
	email_prefix = "%s-parent"
}

resource "squadcast_service" "test" {
	name = "%s"
	description = "some description here."
	team_id = "613611c1eb22db455cfa789f"
	escalation_policy_id = "61361415c2fc70c3101ca7db"
	email_prefix = "%s"
	dependencies = [squadcast_service.test_parent.id]
}

data "squadcast_service" "test" {
	name = squadcast_service.test.name
	team_id = "613611c1eb22db455cfa789f"
}
	`, serviceName, serviceName, serviceName, serviceName)
}
