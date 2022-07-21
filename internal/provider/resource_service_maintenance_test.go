package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
)

func TestAccResourceServiceMaintenance(t *testing.T) {
	resourceName := "squadcast_service_maintenance.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckServiceMaintenanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceServiceMaintenanceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "service_id", "61361611c2fc70c3101ca7dd"),
					resource.TestCheckResourceAttr(resourceName, "windows.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "windows.0.from", "2032-06-01T10:30:00.000Z"),
					resource.TestCheckResourceAttr(resourceName, "windows.0.till", "2032-06-01T11:30:00.000Z"),
					resource.TestCheckResourceAttr(resourceName, "windows.0.repeat_till", ""),
					resource.TestCheckResourceAttr(resourceName, "windows.0.repeat_frequency", ""),
				),
			},
			{
				Config: testAccResourceServiceMaintenanceConfig_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "service_id", "61361611c2fc70c3101ca7dd"),
					resource.TestCheckResourceAttr(resourceName, "windows.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "windows.0.from", "2032-06-01T10:30:00.000Z"),
					resource.TestCheckResourceAttr(resourceName, "windows.0.till", "2032-06-01T11:30:00.000Z"),
					resource.TestCheckResourceAttr(resourceName, "windows.0.repeat_till", "2032-06-30T10:30:00.000Z"),
					resource.TestCheckResourceAttr(resourceName, "windows.0.repeat_frequency", "week"),
					resource.TestCheckResourceAttr(resourceName, "windows.1.from", "2032-07-01T10:30:00.000Z"),
					resource.TestCheckResourceAttr(resourceName, "windows.1.till", "2032-07-02T10:30:00.000Z"),
					resource.TestCheckResourceAttr(resourceName, "windows.1.repeat_till", ""),
					resource.TestCheckResourceAttr(resourceName, "windows.1.repeat_frequency", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "613611c1eb22db455cfa789f:61361611c2fc70c3101ca7dd",
			},
		},
	})
}

func testAccCheckServiceMaintenanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_service_maintenance" {
			continue
		}

		serviceMaintenanceWindows, err := client.GetServiceMaintenanceWindows(context.Background(), rs.Primary.Attributes["service_id"])
		if err != nil {
			return err
		}
		count := len(serviceMaintenanceWindows)
		if count > 0 {
			return fmt.Errorf("expected all service maintenance windows to be destroyed, %d found", count)
		}
	}

	return nil
}

func testAccResourceServiceMaintenanceConfig() string {
	return fmt.Sprintf(`
resource "squadcast_service_maintenance" "test" {
	service_id = "61361611c2fc70c3101ca7dd"

	windows {
		from = "2032-06-01T10:30:00.000Z"
		till = "2032-06-01T11:30:00.000Z"
	}
}
	`)
}

func testAccResourceServiceMaintenanceConfig_update() string {
	return fmt.Sprintf(`
resource "squadcast_service_maintenance" "test" {
	service_id = "61361611c2fc70c3101ca7dd"

	windows {
		from = "2032-06-01T10:30:00.000Z"
		till = "2032-06-01T11:30:00.000Z"
		repeat_till = "2032-06-30T10:30:00.000Z"
		repeat_frequency = "week"
	}

	windows {
		from = "2032-07-01T10:30:00.000Z"
		till = "2032-07-02T10:30:00.000Z"
	}
}
	`)
}
