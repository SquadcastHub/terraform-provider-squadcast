package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
)

func TestAccResourceDelayedNotificationConfig(t *testing.T) {
	fixedResourceName := "squadcast_delayed_notification_config.fixed_test"
	customResourceName := "squadcast_delayed_notification_config.custom_test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckDelayedNotificationConfigDestroy,
		Steps: []resource.TestStep{
			{
				// Test fixed timeslot configuration
				Config: testAccResourceDelayedNotificationConfigFixed(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(fixedResourceName, "id"),
					resource.TestCheckResourceAttr(fixedResourceName, "service_id", "61361611c2fc70c3101ca7dd"),
					resource.TestCheckResourceAttr(fixedResourceName, "is_enabled", "true"),
					resource.TestCheckResourceAttr(fixedResourceName, "timezone", "Asia/Kolkata"),
					resource.TestCheckResourceAttr(fixedResourceName, "fixed_timeslot_config.0.start_time", "09:00"),
					resource.TestCheckResourceAttr(fixedResourceName, "fixed_timeslot_config.0.end_time", "18:00"),
					resource.TestCheckResourceAttr(fixedResourceName, "fixed_timeslot_config.0.repeat_days.#", "3"),
					resource.TestCheckResourceAttr(fixedResourceName, "fixed_timeslot_config.0.repeat_days.0", "sunday"),
					resource.TestCheckResourceAttr(fixedResourceName, "assigned_to.0.id", "61361611c2fc70c3101ca8aa"),
					resource.TestCheckResourceAttr(fixedResourceName, "assigned_to.0.type", "user"),
				),
			},
			{
				// Test import for fixed timeslot config
				ResourceName:      fixedResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "61361611c2fc70c3101ca7dd",
			},
			{
				// Test custom timeslot configuration
				Config: testAccResourceDelayedNotificationConfigCustom(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(customResourceName, "id"),
					resource.TestCheckResourceAttr(customResourceName, "service_id", "61361611c2fc70c3101ca7dd"),
					resource.TestCheckResourceAttr(customResourceName, "is_enabled", "true"),
					resource.TestCheckResourceAttr(customResourceName, "timezone", "Asia/Kolkata"),
					resource.TestCheckResourceAttr(customResourceName, "custom_timeslots_enabled", "true"),
					resource.TestCheckResourceAttr(customResourceName, "custom_timeslots.#", "4"),
					resource.TestCheckResourceAttr(customResourceName, "custom_timeslots.0.day_of_week", "sunday"),
					resource.TestCheckResourceAttr(customResourceName, "custom_timeslots.0.start_time", "10:15"),
					resource.TestCheckResourceAttr(customResourceName, "custom_timeslots.0.end_time", "20:00"),
					resource.TestCheckResourceAttr(customResourceName, "assigned_to.0.id", "61361611c2fc70c3101ca8aa"),
					resource.TestCheckResourceAttr(customResourceName, "assigned_to.0.type", "user"),
				),
			},
			{
				// Test import for custom timeslot config
				ResourceName:      customResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "61361611c2fc70c3101ca7dd",
			},
		},
	})
}

func testAccCheckDelayedNotificationConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_delayed_notification_config" {
			continue
		}

		_, err := client.GetServiceById(context.Background(), "", rs.Primary.Attributes["service_id"])
		if err != nil {
			if api.IsResourceNotFoundError(err) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Delayed Notification Config still exists")
	}

	return nil
}

func testAccResourceDelayedNotificationConfigFixed() string {
	return fmt.Sprintf(`
resource "squadcast_delayed_notification_config" "fixed_test" {
	service_id = "61361611c2fc70c3101ca7dd"
	is_enabled = true
	timezone   = "Asia/Kolkata"

	fixed_timeslot_config {
		start_time  = "09:00"
		end_time    = "18:00"
		repeat_days = ["sunday", "monday", "tuesday"]
	}

	assigned_to {
		id   = "61361611c2fc70c3101ca8aa"
		type = "user"
	}
}
`)
}

func testAccResourceDelayedNotificationConfigCustom() string {
	return fmt.Sprintf(`
resource "squadcast_delayed_notification_config" "custom_test" {
	service_id               = "61361611c2fc70c3101ca7dd"
	is_enabled              = true
	timezone                = "Asia/Kolkata"
	custom_timeslots_enabled = true

	custom_timeslots {
		day_of_week = "sunday"
		start_time  = "10:15"
		end_time    = "20:00"
	}
	custom_timeslots {
		day_of_week = "monday"
		start_time  = "13:15"
		end_time    = "23:59"
	}
	custom_timeslots {
		day_of_week = "tuesday"
		start_time  = "12:15"
		end_time    = "20:59"
	}
	custom_timeslots {
		day_of_week = "wednesday"
		start_time  = "10:15"
		end_time    = "23:59"
	}

	assigned_to {
		id   = "61361611c2fc70c3101ca8aa"
		type = "user"
	}
}
`)
}
