package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
)

func TestAccResourceDeduplicationRules(t *testing.T) {
	resourceName := "squadcast_deduplication_rules.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckDeduplicationRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDeduplicationRulesConfig_defaults(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckTypeSetElemAttr(resourceName, "rules.*", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.is_basic", "false"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.description", "not basic"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.dependency_deduplication", "false"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.time_window", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.time_unit", "hour"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.expression", "payload[\"event_id\"] == 40"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "service_id", "61361611c2fc70c3101ca7dd"),
				),
			},
			{
				Config: testAccResourceDeduplicationRulesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckTypeSetElemAttr(resourceName, "rules.*", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.is_basic", "false"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.description", "not basic"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.dependency_deduplication", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.time_window", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.time_unit", "minute"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.expression", "payload[\"event_id\"] == 40"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "service_id", "61361611c2fc70c3101ca7dd"),
				),
			},
			{
				Config: testAccResourceDeduplicationRulesConfig_updateRules(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckTypeSetElemAttr(resourceName, "rules.*", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.is_basic", "false"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.description", "not basic"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.expression", "payload[\"event_id\"] == 40"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.basic_expressions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.dependency_deduplication", "false"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.time_window", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.time_unit", "hour"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.is_basic", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.description", "basic"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.expression", ""),
					resource.TestCheckResourceAttr(resourceName, "rules.1.basic_expressions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.basic_expressions.0.lhs", "payload[\"foo\"]"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.basic_expressions.0.op", "is"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.basic_expressions.0.rhs", "bar"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.dependency_deduplication", "false"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.time_window", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.time_unit", "hour"),
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

func testAccCheckDeduplicationRulesDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_deduplication_rules" {
			continue
		}

		deduplicationRules, err := client.GetDeduplicationRules(context.Background(), rs.Primary.Attributes["service_id"], rs.Primary.Attributes["team_id"])
		if err != nil {
			return err
		}
		count := len(deduplicationRules.Rules)
		if count > 0 {
			return fmt.Errorf("expected all deduplication rules to be destroyed, %d found", count)
		}
	}

	return nil
}

func testAccResourceDeduplicationRulesConfig_defaults() string {
	return fmt.Sprintf(`
resource "squadcast_deduplication_rules" "test" {
	team_id = "613611c1eb22db455cfa789f"
	service_id = "61361611c2fc70c3101ca7dd"

	rules {
		is_basic = false
		description = "not basic"
		expression = "payload[\"event_id\"] == 40"
	}
}
	`)
}

func testAccResourceDeduplicationRulesConfig() string {
	return fmt.Sprintf(`
resource "squadcast_deduplication_rules" "test" {
	team_id = "613611c1eb22db455cfa789f"
	service_id = "61361611c2fc70c3101ca7dd"

	rules {
		is_basic = false
		description = "not basic"
		expression = "payload[\"event_id\"] == 40"
		dependency_deduplication = true
		time_window = 2
		time_unit = "minute"
	}
}
	`)
}

func testAccResourceDeduplicationRulesConfig_updateRules() string {
	return fmt.Sprintf(`
resource "squadcast_deduplication_rules" "test" {
	team_id = "613611c1eb22db455cfa789f"
	service_id = "61361611c2fc70c3101ca7dd"

	rules {
		is_basic = false
		description = "not basic"
		expression = "payload[\"event_id\"] == 40"
	}

	rules {
		is_basic = true
		description = "basic"

		basic_expressions {
			lhs = "payload[\"foo\"]"
			op = "is"
			rhs = "bar"
		}
	}
}
	`)
}
