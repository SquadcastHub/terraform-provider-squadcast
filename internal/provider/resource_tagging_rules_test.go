package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/testdata"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func TestAccResourceTaggingRules(t *testing.T) {
	teamName := acctest.RandomWithPrefix("test-team")
	user := testdata.RandomUser()
	epName := acctest.RandomWithPrefix("test-ep")
	serviceName := acctest.RandomWithPrefix("test-service")

	teamResourceName := "squadcast_team.test"
	serviceResourceName := "squadcast_service.test"
	resourceName := "squadcast_tagging_rules.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckTaggingRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTaggingRulesConfig(teamName, user, epName, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.is_basic", "false"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.expression", "payload[\"event_id\"] == 40"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.tags.0.key", "MyTag"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.tags.0.value", "foo"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.tags.0.color", "#ababab"),
					resource.TestCheckResourceAttrPair(resourceName, "team_id", teamResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "service_id", serviceResourceName, "id"),
				),
			},
			{
				Config: testAccResourceTaggingRulesConfig_updateRules(teamName, user, epName, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.is_basic", "false"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.expression", "payload[\"event_id\"] == 40"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.basic_expressions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.tags.0.key", "MyTag"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.tags.0.value", "foo"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.tags.0.color", "#ababab"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.is_basic", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.expression", ""),
					resource.TestCheckResourceAttr(resourceName, "rules.1.basic_expressions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.basic_expressions.0.lhs", "payload[\"foo\"]"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.basic_expressions.0.op", "is"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.basic_expressions.0.rhs", "bar"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.tags.0.key", "MyTag"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.tags.0.value", "foo"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.tags.0.color", "#ababab"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.tags.1.key", "MyTag2"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.tags.1.value", "bar"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.tags.1.color", "#f0f0f0"),
					resource.TestCheckResourceAttrPair(resourceName, "team_id", teamResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "service_id", serviceResourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "613611c1eb22db455cfa789f:61361611c2fc70c3101ca7dd",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					teamID, err := tf.StateAttr(s, "squadcast_team", "id")
					if err != nil {
						return "", err
					}

					serviceID, err := tf.StateAttr(s, "squadcast_service", "id")
					if err != nil {
						return "", err
					}

					return teamID + ":" + serviceID, nil
				},
			},
		},
	})
}

func testAccCheckTaggingRulesDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_tagging_rules" {
			continue
		}

		taggingRules, err := client.GetTaggingRules(context.Background(), rs.Primary.Attributes["service_id"], rs.Primary.Attributes["team_id"])
		if err != nil {
			return err
		}
		count := len(taggingRules.Rules)
		if count > 0 {
			return fmt.Errorf("expected all tagging rules to be destroyed, %d found", count)
		}
	}

	return nil
}

func testAccResourceTaggingRulesConfig(teamName string, user testdata.User, epName, serviceName string) string {
	return fmt.Sprintf(`
resource "squadcast_team" "test" {
	name = "%s"
}

resource "squadcast_user" "test" {
	first_name = "%s"
	last_name = "%s"
	email = "%s"
	role = "user"
}

resource "squadcast_team_member" "test" {
	team_id = squadcast_team.test.id
	user_id = squadcast_user.test.id
	role_ids = [
		squadcast_team.test.default_role_ids.admin,
	]
}

resource "squadcast_escalation_policy" "test" {
	name = "%s"

	team_id = squadcast_team.test.id

	rules {
		delay_minutes = 0

		targets {
			id = squadcast_user.test.id
			type = "user"
		}
	}
	depends_on = [squadcast_team_member.test]
}

resource "squadcast_service" "test" {
	name = "%s"
	team_id = squadcast_team.test.id
	escalation_policy_id = squadcast_escalation_policy.test.id
	email_prefix = "testfoo"
}

resource "squadcast_tagging_rules" "test" {
	team_id = squadcast_team.test.id
	service_id = squadcast_service.test.id

	rules {
		is_basic = false
		expression = "payload[\"event_id\"] == 40"

		tags {
			key = "MyTag"
			value = "foo"
			color = "#ababab"
		}
	}
}
	`, teamName, user.FirstName, user.LastName, user.Email, epName, serviceName)
}

func testAccResourceTaggingRulesConfig_updateRules(teamName string, user testdata.User, epName, serviceName string) string {
	return fmt.Sprintf(`
resource "squadcast_team" "test" {
	name = "%s"
}

resource "squadcast_user" "test" {
	first_name = "%s"
	last_name = "%s"
	email = "%s"
	role = "user"
}

resource "squadcast_team_member" "test" {
	team_id = squadcast_team.test.id
	user_id = squadcast_user.test.id
	role_ids = [
		squadcast_team.test.default_role_ids.admin,
	]
}

resource "squadcast_escalation_policy" "test" {
	name = "%s"

	team_id = squadcast_team.test.id

	rules {
		delay_minutes = 0

		targets {
			id = squadcast_user.test.id
			type = "user"
		}
	}
	depends_on = [squadcast_team_member.test]
}

resource "squadcast_service" "test" {
	name = "%s"
	team_id = squadcast_team.test.id
	escalation_policy_id = squadcast_escalation_policy.test.id
	email_prefix = "testfoo"
}

resource "squadcast_tagging_rules" "test" {
	team_id = squadcast_team.test.id
	service_id = squadcast_service.test.id

	rules {
		is_basic = false
		expression = "payload[\"event_id\"] == 40"

		tags {
			key = "MyTag"
			value = "foo"
			color = "#ababab"
		}
	}

	rules {
		is_basic = true

		basic_expressions {
			lhs = "payload[\"foo\"]"
			op = "is"
			rhs = "bar"
		}

		tags {
			key = "MyTag"
			value = "foo"
			color = "#ababab"
		}

		tags {
			key = "MyTag2"
			value = "bar"
			color = "#f0f0f0"
		}
	}
}
	`, teamName, user.FirstName, user.LastName, user.Email, epName, serviceName)
}
