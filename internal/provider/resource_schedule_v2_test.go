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

func TestAccResourceScheduleV2(t *testing.T) {
	scheduleName := acctest.RandomWithPrefix("schedule_v2")

	resourceName := "squadcast_schedule_v2.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckScheduleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScheduleV2Config(scheduleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", scheduleName),
					resource.TestCheckResourceAttr(resourceName, "description", "some description here"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Asia/Kolkata"),
					resource.TestCheckResourceAttr(resourceName, "entity_owner.0.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "entity_owner.0.id", "6113b0ffe4d98ae048c37010"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.key", "key1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", "value1"),
				),
			},
			{
				Config: testAccResourceScheduleV2Config_update(scheduleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", scheduleName),
					resource.TestCheckResourceAttr(resourceName, "description", "some description here"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Asia/Kolkata"),
					resource.TestCheckResourceAttr(resourceName, "entity_owner.0.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "entity_owner.0.id", "6113b0ffe4d98ae048c37010"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.key", "key1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", "value1"),
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

func testAccCheckScheduleV2Destroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_schedule_v2" {
			continue
		}

		_, err := client.GetScheduleV2ById(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("expected schedule to be destroyed, %s found", rs.Primary.ID)
		}

		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceScheduleV2Config(scheduleName string) string {
	return fmt.Sprintf(`
		resource "squadcast_schedule_v2" "test" {
			name = "%s"
			team_id = "613611c1eb22db455cfa789f"
			description = "some description here"
			timezone = "Asia/Kolkata"
			entity_owner {
				type = "user"
				id = "6113b0ffe4d98ae048c37010"
			}
			tags {
				key = "key1"
				value = "value1"
			}
		}
	`, scheduleName)
}

func testAccResourceScheduleV2Config_update(scheduleName string) string {
	return fmt.Sprintf(`
		resource "squadcast_schedule_v2" "test" {
			name = "%s"
			team_id = "613611c1eb22db455cfa789f"
			description = "some description here"
			timezone = "Asia/Kolkata"
			entity_owner {
				type = "user"
				id = "6113b0ffe4d98ae048c37010"
			}
			tags {
				key = "key1"
				value = "value1"
			}
		}
	`, scheduleName)
}
