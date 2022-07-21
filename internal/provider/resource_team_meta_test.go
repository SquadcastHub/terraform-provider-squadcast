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

func TestAccResourceTeam(t *testing.T) {
	teamName := acctest.RandomWithPrefix("team")

	resourceName := "squadcast_team.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTeamConfig(teamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", teamName),
					resource.TestCheckResourceAttr(resourceName, "description", teamName+" description"),
					resource.TestCheckResourceAttr(resourceName, "default", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_role_ids.%", "4"),
					resource.TestCheckResourceAttrSet(resourceName, "default_role_ids.manage_team"),
					resource.TestCheckResourceAttrSet(resourceName, "default_role_ids.admin"),
					resource.TestCheckResourceAttrSet(resourceName, "default_role_ids.user"),
					resource.TestCheckResourceAttrSet(resourceName, "default_role_ids.observer"),
				),
			},
			{
				Config: testAccResourceTeamConfig_update(teamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", teamName+" updated"),
					resource.TestCheckResourceAttr(resourceName, "description", teamName+" description updated"),
					resource.TestCheckResourceAttr(resourceName, "default", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     teamName + " updated",
			},
		},
	})
}

func testAccCheckTeamDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_team" {
			continue
		}

		_, err := client.GetTeamMetaById(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("expected team to be destroyed, %s found", rs.Primary.ID)
		}

		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceTeamConfig(teamName string) string {
	return fmt.Sprintf(`
resource "squadcast_team" "test" {
	name = "%s"
	description = "%s description"
}
	`, teamName, teamName)
}

func testAccResourceTeamConfig_update(teamName string) string {
	return fmt.Sprintf(`
resource "squadcast_team" "test" {
	name = "%s updated"
	description = "%s description updated"
}
	`, teamName, teamName)
}
