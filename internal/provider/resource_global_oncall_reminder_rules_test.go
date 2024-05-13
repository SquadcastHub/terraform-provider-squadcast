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

func TestAccResourceGlobalOncallReminderRules(t *testing.T) {
	teamName := acctest.RandomWithPrefix("test-rules")
	user := testdata.RandomUser()

	teamResourceName := "squadcast_team.test"
	resourceName := "squadcast_global_oncall_reminder_rules.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGlobalOncallReminderRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTeamMemberConfig(teamName, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "team_id", teamResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.type", "Push"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.time", "10"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.type", "Email"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.time", "30"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					teamID, err := tf.StateAttr(s, "squadcast_team", "id")
					if err != nil {
						return "", err
					}

					return teamID + ":" + user.Email, nil
				},
			},
		},
	})
}

func testAccCheckGlobalOncallReminderRulesDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_global_oncall_reminder_rules" {
			continue
		}

		_, err := client.DeleteGlobalOncallReminderRules(context.Background(), rs.Primary.Attributes["team_id"])
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceGlobalOncallReminderRules(teamName string) string {
	return fmt.Sprintf(`
resource "squadcast_team" "test" {
	name = "%s"
}
resource "squadcast_global_oncall_reminder_rules" "test" {
	team_id = squadcast_team.test.id
	rules {
		time = 10
		type = "Push"
	}
	rules {
		time = 30
		type = "Email"
	}
}
	`, teamName)
}
