package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func TestAccResourceTeamRole(t *testing.T) {
	teamRoleName := acctest.RandomWithPrefix("test-teamrole")
	teamName := acctest.RandomWithPrefix("test-team")

	teamResourceName := "squadcast_team.test"
	resourceName := "squadcast_team_role.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckTeamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTeamRoleConfig(teamName, teamRoleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "team_id", teamResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", teamRoleName),
					resource.TestCheckResourceAttr(resourceName, "default", "false"),
					resource.TestCheckResourceAttr(resourceName, "abilities.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "abilities.0", "read-escalation-policies"),
				),
			},
			{
				Config: testAccResourceTeamRoleConfig_update(teamName, teamRoleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "team_id", teamResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", teamRoleName),
					resource.TestCheckResourceAttr(resourceName, "default", "false"),
					resource.TestCheckResourceAttr(resourceName, "abilities.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "abilities.0", "read-escalation-policies"),
					resource.TestCheckResourceAttr(resourceName, "abilities.1", "update-runbooks"),
				),
			},
			{
				Config: testAccResourceTeamRoleConfig_noAbilities(teamName, teamRoleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "team_id", teamResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", teamRoleName),
					resource.TestCheckResourceAttr(resourceName, "default", "false"),
					resource.TestCheckResourceAttr(resourceName, "abilities.#", "0"),
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

					return teamID + ":" + teamRoleName, nil
				},
			},
		},
	})
}

func testAccCheckTeamRoleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_team_role" {
			continue
		}

		_, err := client.GetTeamRoleByID(context.Background(), rs.Primary.Attributes["team_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("expected team role to be destroyed, %s found", rs.Primary.ID)
		}

		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceTeamRoleConfig(teamName, teamRoleName string) string {
	return fmt.Sprintf(`
resource "squadcast_team" "test" {
	name = "%s"
}

resource "squadcast_team_role" "test" {
	name = "%s"
	team_id = squadcast_team.test.id
	abilities = ["read-escalation-policies"]
}
	`, teamName, teamRoleName)
}

func testAccResourceTeamRoleConfig_update(teamName, teamRoleName string) string {
	return fmt.Sprintf(`
resource "squadcast_team" "test" {
	name = "%s"
}

resource "squadcast_team_role" "test" {
	name = "%s"
	team_id = squadcast_team.test.id
	abilities = ["read-escalation-policies", "update-runbooks"]
}
	`, teamName, teamRoleName)
}

func testAccResourceTeamRoleConfig_noAbilities(teamName, teamRoleName string) string {
	return fmt.Sprintf(`
resource "squadcast_team" "test" {
	name = "%s"
}

resource "squadcast_team_role" "test" {
	name = "%s"
	team_id = squadcast_team.test.id
	abilities = []
}
	`, teamName, teamRoleName)
}
