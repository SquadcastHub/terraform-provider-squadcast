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

func TestAccResourceTeamMember(t *testing.T) {
	teamName := acctest.RandomWithPrefix("test-team")
	user := testdata.RandomUser()

	teamResourceName := "squadcast_team.test"
	userResourceName := "squadcast_user.test"
	resourceName := "squadcast_team_member.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckTeamMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTeamMemberConfig(teamName, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "team_id", teamResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "user_id", userResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "role_ids.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "role_ids.0", teamResourceName, "default_role_ids.admin"),
					resource.TestCheckResourceAttrPair(resourceName, "role_ids.1", teamResourceName, "default_role_ids.user"),
				),
			},
			{
				Config: testAccResourceTeamMemberConfig_observer(teamName, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "team_id", teamResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "user_id", userResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "role_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "role_ids.0", teamResourceName, "default_role_ids.observer"),
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

func testAccCheckTeamMemberDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_team_member" {
			continue
		}

		_, err := client.GetTeamMemberByID(context.Background(), rs.Primary.Attributes["team_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("expected member to be deleted, but was found")
		}

		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceTeamMemberConfig(teamName string, user testdata.User) string {
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
		squadcast_team.test.default_role_ids.user
	]
}
	`, teamName, user.FirstName, user.LastName, user.Email)
}

func testAccResourceTeamMemberConfig_observer(teamName string, user testdata.User) string {
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
			squadcast_team.test.default_role_ids.observer,
		]
	}
	`, teamName, user.FirstName, user.LastName, user.Email)
}
