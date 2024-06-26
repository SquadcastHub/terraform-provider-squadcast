package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceEscalationPolicy(t *testing.T) {
	escalationPolicyName := acctest.RandomWithPrefix("escalation_policy")

	resourceName := "data.squadcast_escalation_policy.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEscalationPolicyConfig(escalationPolicyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", escalationPolicyName),
					resource.TestCheckResourceAttr(resourceName, "description", "It's an amazing policy"),
					resource.TestCheckResourceAttr(resourceName, "repeat.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "repeat.0.times", "2"),
					resource.TestCheckResourceAttr(resourceName, "repeat.0.delay_minutes", "10"),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.delay_minutes", "0"),
					resource.TestCheckNoResourceAttr(resourceName, "rules.0.notification_channels.#"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.targets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.targets.0.id", "5f8891527f735f0a6646f3b7"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.targets.0.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.targets.1.id", "5eb26b36ec9f070550204c85"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.targets.1.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "entity_owner.id", "6113b0ffe4d98ae048c37010"),
					resource.TestCheckResourceAttr(resourceName, "entity_owner.type", "user"),
				),
			},
			{
				Config: testAccDataSourceEscalationPolicyConfig_2rules(escalationPolicyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", escalationPolicyName),
					resource.TestCheckResourceAttr(resourceName, "description", "It's an amazing policy"),
					resource.TestCheckResourceAttr(resourceName, "entity_owner.id", "6113b0ffe4d98ae048c37010"),
					resource.TestCheckResourceAttr(resourceName, "entity_owner.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "repeat.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "repeat.0.times", "2"),
					resource.TestCheckResourceAttr(resourceName, "repeat.0.delay_minutes", "10"),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.delay_minutes", "0"),
					resource.TestCheckNoResourceAttr(resourceName, "rules.0.notification_channels.#"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.targets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.targets.0.id", "5f8891527f735f0a6646f3b7"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.targets.0.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.targets.1.id", "5eb26b36ec9f070550204c85"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.targets.1.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.delay_minutes", "5"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.notification_channels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.notification_channels.0", "Phone"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.targets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.targets.0.id", "61c98f3c75b3a4ebc787f88e"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.targets.0.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.targets.1.id", "5ef5de4259c32c7ca25b0bfa"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.targets.1.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.repeat.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.repeat.0.times", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.repeat.0.delay_minutes", "5"),
				),
			},
			{
				Config: testAccDataSourceEscalationPolicyConfig_3rules(escalationPolicyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", escalationPolicyName),
					resource.TestCheckResourceAttr(resourceName, "description", "It's an amazing policy"),
					resource.TestCheckResourceAttr(resourceName, "entity_owner.id", "6113b0ffe4d98ae048c37010"),
					resource.TestCheckResourceAttr(resourceName, "entity_owner.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "repeat.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "repeat.0.times", "2"),
					resource.TestCheckResourceAttr(resourceName, "repeat.0.delay_minutes", "10"),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.delay_minutes", "0"),
					resource.TestCheckNoResourceAttr(resourceName, "rules.0.notification_channels.#"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.targets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.targets.0.id", "5f8891527f735f0a6646f3b7"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.targets.0.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.targets.1.id", "5eb26b36ec9f070550204c85"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.targets.1.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.delay_minutes", "5"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.notification_channels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.notification_channels.0", "Phone"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.targets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.targets.0.id", "61c98f3c75b3a4ebc787f88e"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.targets.0.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.targets.1.id", "5ef5de4259c32c7ca25b0bfa"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.targets.1.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.repeat.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.repeat.0.times", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.repeat.0.delay_minutes", "5"),
					resource.TestCheckResourceAttr(resourceName, "rules.2.delay_minutes", "10"),
					resource.TestCheckNoResourceAttr(resourceName, "rules.2.notification_channels.#"),
					resource.TestCheckResourceAttr(resourceName, "rules.2.targets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.2.targets.0.id", "60b8bcd7ff5010bf96583e03"),
					resource.TestCheckResourceAttr(resourceName, "rules.2.targets.0.type", "squad"),
					resource.TestCheckResourceAttr(resourceName, "rules.2.targets.1.id", "62a6242c40977285b03b57e3"),
					resource.TestCheckResourceAttr(resourceName, "rules.2.targets.1.type", "schedule"),
					resource.TestCheckResourceAttr(resourceName, "rules.2.round_robin.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.2.round_robin.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.2.round_robin.0.rotation.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.2.round_robin.0.rotation.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.2.round_robin.0.rotation.0.delay_minutes", "1"),
					resource.TestCheckNoResourceAttr(resourceName, "rules.2.repeat.#"),
				),
			},
		},
	})
}

func testAccDataSourceEscalationPolicyConfig(escalationPolicyName string) string {
	return fmt.Sprintf(`
resource "squadcast_escalation_policy" "test" {
	name = "%s"
	description = "It's an amazing policy"

	team_id = "613611c1eb22db455cfa789f"

	entity_owner{
		id = "6113b0ffe4d98ae048c37010"
		type = "user"
	}

	rules {
		delay_minutes = 0

		targets {
			id = "5f8891527f735f0a6646f3b7"
			type = "user"
		}

		targets {
			id = "5eb26b36ec9f070550204c85"
			type = "user"
		}
	}

	repeat {
        times = 2
        delay_minutes = 10
    }
}

data "squadcast_escalation_policy" "test" {
	name = squadcast_escalation_policy.test.name
	team_id = "613611c1eb22db455cfa789f"
}
	`, escalationPolicyName)
}

func testAccDataSourceEscalationPolicyConfig_2rules(escalationPolicyName string) string {
	return fmt.Sprintf(`
resource "squadcast_escalation_policy" "test" {
	name = "%s"
	description = "It's an amazing policy"

	team_id = "613611c1eb22db455cfa789f"

	rules {
		delay_minutes = 0

		targets {
			id = "5f8891527f735f0a6646f3b7"
			type = "user"
		}

		targets {
			id = "5eb26b36ec9f070550204c85"
			type = "user"
		}
	}

	rules {
        delay_minutes = 5

        targets {
            id = "61c98f3c75b3a4ebc787f88e"
            type = "user"
        }

        targets {
            id = "5ef5de4259c32c7ca25b0bfa"
            type = "user"
        }

        notification_channels = ["Phone"]

        repeat {
            times = 1
            delay_minutes = 5
        }
    }

	repeat {
        times = 2
        delay_minutes = 10
    }
}

data "squadcast_escalation_policy" "test" {
	name = squadcast_escalation_policy.test.name
	team_id = "613611c1eb22db455cfa789f"
}
	`, escalationPolicyName)
}

func testAccDataSourceEscalationPolicyConfig_3rules(escalationPolicyName string) string {
	return fmt.Sprintf(`
resource "squadcast_escalation_policy" "test" {
	name = "%s"
	description = "It's an amazing policy"

	team_id = "613611c1eb22db455cfa789f"

	rules {
		delay_minutes = 0

		targets {
			id = "5f8891527f735f0a6646f3b7"
			type = "user"
		}

		targets {
			id = "5eb26b36ec9f070550204c85"
			type = "user"
		}
	}

	rules {
        delay_minutes = 5

        targets {
            id = "61c98f3c75b3a4ebc787f88e"
            type = "user"
        }

        targets {
            id = "5ef5de4259c32c7ca25b0bfa"
            type = "user"
        }

        notification_channels = ["Phone"]

        repeat {
            times = 1
            delay_minutes = 5
        }
    }

	rules {
        delay_minutes = 10

        targets {
            id = "60b8bcd7ff5010bf96583e03"
            type = "squad"
        }

        targets {
            id = "62a6242c40977285b03b57e3"
            type = "schedule"
        }

        round_robin {
            enabled = true

            rotation {
                enabled = true
                delay_minutes = 1
            }
        }
    }

	repeat {
        times = 2
        delay_minutes = 10
    }
}

data "squadcast_escalation_policy" "test" {
	name = squadcast_escalation_policy.test.name
	team_id = "613611c1eb22db455cfa789f"
}
	`, escalationPolicyName)
}
