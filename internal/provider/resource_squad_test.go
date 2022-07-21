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

func TestAccResourceSquad(t *testing.T) {
	squadName := acctest.RandomWithPrefix("squad")

	resourceName := "squadcast_squad.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckSquadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSquadConfig(squadName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckTypeSetElemAttr(resourceName, "member_ids.*", "1"),
					resource.TestCheckResourceAttr(resourceName, "member_ids.0", "5f8891527f735f0a6646f3b6"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", squadName),
				),
			},
			{
				Config: testAccResourceSquadConfig_updateMembers(squadName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckTypeSetElemAttr(resourceName, "member_ids.*", "2"),
					resource.TestCheckResourceAttr(resourceName, "member_ids.0", "5f8891527f735f0a6646f3b6"),
					resource.TestCheckResourceAttr(resourceName, "member_ids.1", "5eb26b36ec9f070550204c85"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", squadName),
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

func testAccCheckSquadDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_squad" {
			continue
		}

		_, err := client.GetSquadById(context.Background(), rs.Primary.Attributes["team_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("expected squad to be destroyed, %s found", rs.Primary.ID)
		}

		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceSquadConfig(squadName string) string {
	return fmt.Sprintf(`
resource "squadcast_squad" "test" {
	name = "%s"
	team_id = "613611c1eb22db455cfa789f"
	member_ids = ["5f8891527f735f0a6646f3b6"]
}
	`, squadName)
}

func testAccResourceSquadConfig_updateMembers(squadName string) string {
	return fmt.Sprintf(`
resource "squadcast_squad" "test" {
	name = "%s"
	team_id = "613611c1eb22db455cfa789f"
	member_ids = ["5f8891527f735f0a6646f3b6", "5eb26b36ec9f070550204c85"]
}
	`, squadName)
}
