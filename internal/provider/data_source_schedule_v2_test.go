package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceScheduleV2(t *testing.T) {
	scheduleName := acctest.RandomWithPrefix("schedule_v2")

	resourceName := "data.squadcast_schedule_v2.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccScheduleV2DataSourceConfig(scheduleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", scheduleName),
					resource.TestCheckResourceAttr(resourceName, "description", "Test schedule"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Asia/Kolkata"),
					resource.TestCheckResourceAttr(resourceName, "entity_owner.0.type", "team"),
					resource.TestCheckResourceAttr(resourceName, "entity_owner.0.id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.key", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.color", "#9900ef"),
				),
			},
		},
	})
}

func testAccScheduleV2DataSourceConfig(scheduleName string) string {
	return fmt.Sprintf(`
		resource "squadcast_schedule_v2" "test" {
			name = "%s"
			team_id = "613611c1eb22db455cfa789f"
			timezone = "Asia/Kolkata"
			description = "Test schedule"
			entity_owner {
				type = "team"
				id = "613611c1eb22db455cfa789f"
			}
			tags {
				key = "test"
				value = "test"
				color = "#9900ef"
			}
		}

		data "squadcast_schedule_v2" "test" {
			name = squadcast_schedule_v2.test.name
			team_id = "613611c1eb22db455cfa789f"
		}
	`, scheduleName)
}
