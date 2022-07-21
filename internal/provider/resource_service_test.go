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

func TestAccResourceService(t *testing.T) {
	serviceName := acctest.RandomWithPrefix("service")

	resourceName := "squadcast_service.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceServiceConfig(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", serviceName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "escalation_policy_id", "5f8c4ff09b0ccd917237c04b"),
					resource.TestCheckResourceAttr(resourceName, "email_prefix", "testfoo"),
					resource.TestCheckResourceAttrSet(resourceName, "api_key"),
					resource.TestCheckResourceAttr(resourceName, "email", "testfoo@squadcast.incidents.squadcast.com"),
					resource.TestCheckResourceAttr(resourceName, "dependencies.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "alert_source_endpoints.email", "testfoo@squadcast.incidents.squadcast.com"),
				),
			},
			{
				Config: testAccResourceServiceConfig_update(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", serviceName),
					resource.TestCheckResourceAttr(resourceName, "description", "some description here."),
					resource.TestCheckResourceAttr(resourceName, "escalation_policy_id", "61361415c2fc70c3101ca7db"),
					resource.TestCheckResourceAttr(resourceName, "email_prefix", "foomp2"),
					resource.TestCheckResourceAttrSet(resourceName, "api_key"),
					resource.TestCheckResourceAttr(resourceName, "email", "foomp2@squadcast.incidents.squadcast.com"),
					resource.TestCheckResourceAttr(resourceName, "dependencies.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "alert_source_endpoints.email", "foomp2@squadcast.incidents.squadcast.com"),
				),
			},
			{
				Config: testAccResourceServiceConfig_dependencies(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", serviceName),
					resource.TestCheckResourceAttr(resourceName, "description", "some description here."),
					resource.TestCheckResourceAttr(resourceName, "escalation_policy_id", "61361415c2fc70c3101ca7db"),
					resource.TestCheckResourceAttr(resourceName, "email_prefix", "foomp2"),
					resource.TestCheckResourceAttrSet(resourceName, "api_key"),
					resource.TestCheckResourceAttr(resourceName, "email", "foomp2@squadcast.incidents.squadcast.com"),
					resource.TestCheckResourceAttr(resourceName, "dependencies.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "dependencies.0", "squadcast_service.test_parent", "id"),
					resource.TestCheckResourceAttr(resourceName, "alert_source_endpoints.email", "foomp2@squadcast.incidents.squadcast.com"),
				),
			},
			{
				ResourceName:        resourceName,
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: "613611c1eb22db455cfa789f:",
			},
		},
	})
}

func testAccCheckServiceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_service" {
			continue
		}

		_, err := client.GetServiceById(context.Background(), rs.Primary.Attributes["team_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("expected service to be destroyed, %s found", rs.Primary.ID)
		}

		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceServiceConfig(serviceName string) string {
	return fmt.Sprintf(`
resource "squadcast_service" "test" {
	name = "%s"
	team_id = "613611c1eb22db455cfa789f"
	escalation_policy_id = "5f8c4ff09b0ccd917237c04b"
	email_prefix = "testfoo"
}
	`, serviceName)
}

func testAccResourceServiceConfig_update(serviceName string) string {
	return fmt.Sprintf(`
resource "squadcast_service" "test" {
	name = "%s"
	description = "some description here."
	team_id = "613611c1eb22db455cfa789f"
	escalation_policy_id = "61361415c2fc70c3101ca7db"
	email_prefix = "foomp2"
}
	`, serviceName)
}

func testAccResourceServiceConfig_dependencies(serviceName string) string {
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
	email_prefix = "foomp2"
	dependencies = [squadcast_service.test_parent.id]
}
	`, serviceName, serviceName, serviceName)
}
