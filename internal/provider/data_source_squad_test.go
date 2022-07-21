package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSquad(t *testing.T) {
	squadName := acctest.RandomWithPrefix("squad")

	resourceName := "data.squadcast_squad.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSquadDataSourceConfig(squadName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckTypeSetElemAttr(resourceName, "member_ids.*", "1"),
					resource.TestCheckResourceAttr(resourceName, "member_ids.0", "5f8891527f735f0a6646f3b6"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", squadName),
				),
			},
		},
	})
}

func testAccSquadDataSourceConfig(squadName string) string {
	return fmt.Sprintf(`
resource "squadcast_squad" "test" {
	name = "%s"
	team_id = "613611c1eb22db455cfa789f"
	member_ids = ["5f8891527f735f0a6646f3b6"]
}

data "squadcast_squad" "test" {
	name = squadcast_squad.test.name
	team_id = "613611c1eb22db455cfa789f"
}
	`, squadName)
}
