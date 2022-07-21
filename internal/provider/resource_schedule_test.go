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

func TestAccResourceSchedule(t *testing.T) {
	scheduleName := acctest.RandomWithPrefix("schedule")

	resourceName := "squadcast_schedule.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScheduleConfig(scheduleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", scheduleName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "color", "#9900ef"),
				),
			},
			{
				Config: testAccResourceScheduleConfig_update(scheduleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", scheduleName),
					resource.TestCheckResourceAttr(resourceName, "description", "some description here"),
					resource.TestCheckResourceAttr(resourceName, "color", "#fff000"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "613611c1eb22db455cfa789f:" + scheduleName,
			},
		},
	})
}

func testAccCheckScheduleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_schedule" {
			continue
		}

		_, err := client.GetScheduleById(context.Background(), rs.Primary.Attributes["team_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("expected schedule to be destroyed, %s found", rs.Primary.ID)
		}

		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceScheduleConfig(scheduleName string) string {
	return fmt.Sprintf(`
resource "squadcast_schedule" "test" {
	name = "%s"
	team_id = "613611c1eb22db455cfa789f"
	color = "#9900ef"
}
	`, scheduleName)
}

func testAccResourceScheduleConfig_update(scheduleName string) string {
	return fmt.Sprintf(`
resource "squadcast_schedule" "test" {
	name = "%s"
	team_id = "613611c1eb22db455cfa789f"
	description = "some description here"
	color = "#fff000"
}
	`, scheduleName)
}
