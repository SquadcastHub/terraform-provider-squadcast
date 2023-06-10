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
		// PreCheck:          func() { testAccPreCheck(t) },
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
					resource.TestCheckResourceAttr(resourceName, "entity_owner.0.type", "team"),
					resource.TestCheckResourceAttr(resourceName, "entity_owner.0.id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.key", "key1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.color", "#9900ef"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.name", "Test Rotation"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.start_date", "2023-06-09T00:00:00Z"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.period", "custom"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.change_participants_frequency", "1"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.change_participants_unit", "week"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.custom_period_frequency", "1"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.custom_period_unit", "week"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.participant_groups.0.participants.0.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.participant_groups.0.participants.0.id", "613611c1eb22db455cfa789f"),
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
					resource.TestCheckResourceAttr(resourceName, "entity_owner.0.type", "team"),
					resource.TestCheckResourceAttr(resourceName, "entity_owner.0.id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.key", "key1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.color", "#9900ef"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.name", "Test Rotation"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.start_date", "2023-06-09T00:00:00Z"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.period", "custom"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.change_participants_frequency", "1"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.change_participants_unit", "week"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.custom_period_frequency", "1"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.custom_period_unit", "week"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.participant_groups.0.participants.0.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "rotation.0.participant_groups.0.participants.0.id", "613611c1eb22db455cfa789f"),
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
				type = "team"
				id = "613611c1eb22db455cfa789f"
			}
			tags {
				key = "key1"
				value = "value1"
				color = "#9900ef"
			}
			rotation {
				name = "Test Rotation"
				start_date = "2023-06-09T00:00:00Z"
				period = "custom"
				change_participants_frequency = 1
				change_participants_unit = "week"
				custom_period_frequency = 1
				custom_period_unit = "week"
				participant_groups {
					participants {
						type = "user"
						id = "613611c1eb22db455cfa789f"
					}
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
				type = "team"
				id = "613611c1eb22db455cfa789f"
			}
			tags {
				key = "key1"
				value = "value1"
				color = "#9900ef"
			}
			rotation {
				name = "Test Rotation"
				start_date = "2023-06-09T00:00:00Z"
				period = "custom"
				change_participants_frequency = 1
				change_participants_unit = "week"
				custom_period_frequency = 1
				custom_period_unit = "week"
				participant_groups {
					participants {
						type = "user"
						id = "613611c1eb22db455cfa789f"
					}
				}
		}
	`, scheduleName)
}
