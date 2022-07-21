package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTeam(t *testing.T) {
	resourceName := "data.squadcast_team.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "Default Team"),
					resource.TestCheckResourceAttr(resourceName, "description", "Default team"),
					resource.TestCheckResourceAttr(resourceName, "default", "true"),
					resource.TestCheckResourceAttr(resourceName, "members.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "members.0.user_id", "5ef5de4259c32c7ca25b0bfa"),
					resource.TestCheckResourceAttr(resourceName, "members.0.role_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.0.role_ids.0", "613611c1eb22db455cfa789b"),
					resource.TestCheckResourceAttr(resourceName, "members.0.role_ids.1", "613611c1eb22db455cfa789c"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "roles.0.id", "613611c1eb22db455cfa789b"),
					resource.TestCheckResourceAttr(resourceName, "roles.0.name", "Manage Team"),
					resource.TestCheckResourceAttr(resourceName, "roles.0.default", "true"),
					resource.TestCheckResourceAttr(resourceName, "roles.0.abilities.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "roles.0.abilities.0", "delete-teams"),
					resource.TestCheckResourceAttr(resourceName, "roles.0.abilities.1", "read-teams"),
					resource.TestCheckResourceAttr(resourceName, "roles.0.abilities.2", "update-teams"),
				),
			},
		},
	})
}

func testAccTeamDataSourceConfig() string {
	return fmt.Sprintf(`
data "squadcast_team" "test" {
	name = "Default Team"
}
	`)
}
