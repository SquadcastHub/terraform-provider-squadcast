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

func TestAccResourceScheduleRotation(t *testing.T) {
	rotationName := acctest.RandomWithPrefix("schedule_rotation_v2")

	resourceName := "squadcast_schedule_rotation_v2.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckScheduleRotationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScheduleRotationConfig(rotationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "schedule_id", "100"),
					resource.TestCheckResourceAttr(resourceName, "name", rotationName),
					resource.TestCheckResourceAttr(resourceName, "start_date", "2023-06-09T00:00:00Z"),
					resource.TestCheckResourceAttr(resourceName, "period", "weekly"),
					resource.TestCheckResourceAttr(resourceName, "shift_timeslots.0.start_hour", "10"),
					resource.TestCheckResourceAttr(resourceName, "shift_timeslots.0.start_minute", "30"),
					resource.TestCheckResourceAttr(resourceName, "shift_timeslots.0.duration", "720"),
					resource.TestCheckResourceAttr(resourceName, "change_participants_frequency", "1"),
					resource.TestCheckResourceAttr(resourceName, "change_participants_unit", "rotation"),
					resource.TestCheckResourceAttr(resourceName, "participant_groups.0.participants.0.type", "team"),
					resource.TestCheckResourceAttr(resourceName, "participant_groups.0.participants.0.id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "ends_after_iterations", "2"),
				),
			},
			{
				Config: testAccResourceScheduleRotationConfig_update(rotationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "schedule_id", "100"),
					resource.TestCheckResourceAttr(resourceName, "name", rotationName),
					resource.TestCheckResourceAttr(resourceName, "start_date", "2023-06-13T00:00:00Z"),
					resource.TestCheckResourceAttr(resourceName, "period", "custom"),
					resource.TestCheckResourceAttr(resourceName, "shift_timeslots.0.start_hour", "10"),
					resource.TestCheckResourceAttr(resourceName, "shift_timeslots.0.start_minute", "30"),
					resource.TestCheckResourceAttr(resourceName, "shift_timeslots.0.duration", "1440"),
					resource.TestCheckResourceAttr(resourceName, "shift_timeslots.0.day_of_week", "saturday"),
					resource.TestCheckResourceAttr(resourceName, "shift_timeslots.1.start_hour", "12"),
					resource.TestCheckResourceAttr(resourceName, "shift_timeslots.1.start_minute", "30"),
					resource.TestCheckResourceAttr(resourceName, "shift_timeslots.1.duration", "720"),
					resource.TestCheckResourceAttr(resourceName, "shift_timeslots.1.day_of_week", "sunday"),
					resource.TestCheckResourceAttr(resourceName, "change_participants_frequency", "1"),
					resource.TestCheckResourceAttr(resourceName, "change_participants_unit", "rotation"),
					resource.TestCheckResourceAttr(resourceName, "custom_period_frequency", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_period_unit", "week"),
					resource.TestCheckResourceAttr(resourceName, "participant_groups.0.participants.0.type", "team"),
					resource.TestCheckResourceAttr(resourceName, "participant_groups.0.participants.0.id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "end_date", "2023-08-31T00:00:00Z"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "613611c1eb22db455cfa789f:Test Schedule:" + rotationName,
			},
		},
	})
}

func testAccCheckScheduleRotationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_schedule_rotation_v2" {
			continue
		}

		_, err := client.GetScheduleRotationById(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("expected rotation to be destroyed, %s found", rs.Primary.ID)
		}

		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceScheduleRotationConfig(rotationName string) string {
	return fmt.Sprintf(`
		resource "squadcast_schedule_rotation_v2" "test" {
			schedule_id = "100"
			name = "%s"
			start_date = "2023-07-01T00:00:00Z"
			period = "weekly"
			shift_timeslots {
				start_hour = 10
				start_minute = 30
				duration = 720
			}
			change_participants_frequency = 1
			change_participants_unit = "rotation"
			participant_groups {
				participants {
					id = "613611c1eb22db455cfa789f"
					type = "team"
				}
			}
			ends_after_iterations = 2
		}
	`, rotationName)
}

func testAccResourceScheduleRotationConfig_update(rotationName string) string {
	return fmt.Sprintf(`
		resource "squadcast_schedule_rotation_v2" "test" {
			schedule_id = "100"
			name = "%s"
			start_date = "2023-06-13T00:00:00Z"
			period = "custom"
			shift_timeslots {
				start_hour = 10
				start_minute = 0
				duration = 1440
				day_of_week = "saturday"
			}
			shift_timeslots {
				start_hour = 12
				start_minute = 30
				duration = 720
				day_of_week = "sunday"
			}
			change_participants_frequency = 1
			change_participants_unit = "rotation"
			custom_period_frequency = 1
			custom_period_unit = "week"
			participant_groups {
				participants {
					id = "613611c1eb22db455cfa789f"
					type = "team"
				}
			}
			end_date =  "2023-08-31T00:00:00Z"
		}
	`, rotationName)
}
