package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSchedule(t *testing.T) {
	scheduleName := acctest.RandomWithPrefix("schedule")

	resourceName := "data.squadcast_schedule.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccScheduleDataSourceConfig(scheduleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", scheduleName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "color", "#9900ef"),
				),
			},
		},
	})
}

func testAccScheduleDataSourceConfig(scheduleName string) string {
	return fmt.Sprintf(`
resource "squadcast_schedule" "test" {
	name = "%s"
	team_id = "613611c1eb22db455cfa789f"
	color = "#9900ef"
}

data "squadcast_schedule" "test" {
	name = squadcast_schedule.test.name
	team_id = "613611c1eb22db455cfa789f"
}
	`, scheduleName)
}
