package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
)

func TestAccResourceRoutingRules(t *testing.T) {
	resourceName := "squadcast_routing_rules.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckRoutingRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRoutingRulesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckTypeSetElemAttr(resourceName, "rules.*", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.is_basic", "false"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.route_to_id", "5f8891527f735f0a6646f3b6"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.route_to_type", "user"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.expression", "payload[\"event_id\"] == 40"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "service_id", "61361611c2fc70c3101ca7dd"),
				),
			},
			{
				Config: testAccResourceRoutingRulesConfig_updateRules(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckTypeSetElemAttr(resourceName, "rules.*", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.is_basic", "false"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.expression", "payload[\"event_id\"] == 40"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.basic_expressions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.route_to_id", "5f8891527f735f0a6646f3b6"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.route_to_type", "user"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.is_basic", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.expression", ""),
					resource.TestCheckResourceAttr(resourceName, "rules.1.basic_expressions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.basic_expressions.0.lhs", "payload[\"foo\"]"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.basic_expressions.0.rhs", "bar"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.route_to_id", "5f8c4ff09b0ccd917237c04b"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.route_to_type", "escalationpolicy"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "service_id", "61361611c2fc70c3101ca7dd"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "613611c1eb22db455cfa789f:61361611c2fc70c3101ca7dd",
			},
		},
	})
}

func testAccCheckRoutingRulesDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_routing_rules" {
			continue
		}

		routingRules, err := client.GetRoutingRules(context.Background(), rs.Primary.Attributes["service_id"], rs.Primary.Attributes["team_id"])
		if err != nil {
			return err
		}
		count := len(routingRules.Rules)
		if count > 0 {
			return fmt.Errorf("expected all routing rules to be destroyed, %d found", count)
		}
	}

	return nil
}

func testAccResourceRoutingRulesConfig() string {
	return fmt.Sprintf(`
resource "squadcast_routing_rules" "test" {
	team_id = "613611c1eb22db455cfa789f"
	service_id = "61361611c2fc70c3101ca7dd"

	rules {
		is_basic = false
		expression = "payload[\"event_id\"] == 40"
		route_to_id = "5f8891527f735f0a6646f3b6"
		route_to_type = "user"
	}
}
	`)
}

func testAccResourceRoutingRulesConfig_updateRules() string {
	return fmt.Sprintf(`
resource "squadcast_routing_rules" "test" {
	team_id = "613611c1eb22db455cfa789f"
	service_id = "61361611c2fc70c3101ca7dd"

	rules {
		is_basic = false
		expression = "payload[\"event_id\"] == 40"

		route_to_id = "5f8891527f735f0a6646f3b6"
		route_to_type = "user"
	}

	rules {
		is_basic = true

		basic_expressions {
			lhs = "payload[\"foo\"]"
			rhs = "bar"
		}

		route_to_id = "5f8c4ff09b0ccd917237c04b"
		route_to_type = "escalationpolicy"
	}
}
	`)
}
