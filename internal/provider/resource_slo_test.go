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

func TestAccResourceSlo(t *testing.T) {
	sloName := acctest.RandomWithPrefix("terraform-acc-test-slo-")

	resourceName := "squadcast_slo.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckSloDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSloConfig(sloName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", sloName),
					resource.TestCheckResourceAttr(resourceName, "description", "Tracks some slo for some service"),
					resource.TestCheckResourceAttr(resourceName, "target_slo", "99.9"),
					resource.TestCheckResourceAttr(resourceName, "duration_in_days", "30"),
					resource.TestCheckResourceAttr(resourceName, "time_interval_type", "rolling"),
					resource.TestCheckResourceAttr(resourceName, "service_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "service_ids.0", "6257a8eb3c8ff45615ce5f2e"),
					resource.TestCheckResourceAttr(resourceName, "slis.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "slis.0", "latency"),
					resource.TestCheckResourceAttr(resourceName, "notify.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notify.0.user_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notify.0.user_ids.0", "6113b0ffe4d98ae048c37010"),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.name", "breached_error_budget"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.name", "unhealthy_slo"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.threshold", "1"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "61443b953ffd52818bf1616a"),
				),
			},
			{
				Config: testAccResourceSloConfig_update(sloName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", sloName),
					resource.TestCheckResourceAttr(resourceName, "description", "Tracks some slo for some test service"),
					resource.TestCheckResourceAttr(resourceName, "target_slo", "99.99"),
					resource.TestCheckResourceAttr(resourceName, "duration_in_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "time_interval_type", "rolling"),
					resource.TestCheckResourceAttr(resourceName, "service_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "service_ids.0", "6257a8eb3c8ff45615ce5f2e"),
					resource.TestCheckResourceAttr(resourceName, "slis.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "slis.0", "latency"),
					resource.TestCheckResourceAttr(resourceName, "slis.1", "high-err-rate"),
					resource.TestCheckResourceAttr(resourceName, "notify.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notify.0.user_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "notify.0.user_ids.0", "6113b0ffe4d98ae048c37010"),
					resource.TestCheckResourceAttr(resourceName, "notify.0.user_ids.1", "61305a78127c63c6d2c8f746"),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.name", "breached_error_budget"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.name", "unhealthy_slo"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.threshold", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.2.name", "remaining_error_budget"),
					resource.TestCheckResourceAttr(resourceName, "rules.2.threshold", "11"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "61443b953ffd52818bf1616a"),
				),
			},
		},
	})
}

func testAccCheckSloDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_slo" {
			continue
		}

		slo, _ := client.GetSlo(context.Background(), client.OrganizationID, rs.Primary.Attributes["id"], rs.Primary.Attributes["team_id"])
		if slo != nil {
			return fmt.Errorf("expected slo to be destroyed, %s found", slo.Name)
		}
	}
	return nil
}

func testAccResourceSloConfig(sloName string) string {
	return fmt.Sprintf(`

resource "squadcast_slo" "test" {
	name = "%s"
	description = "Tracks some slo for some service"
	target_slo = 99.9
	service_ids = ["6257a8eb3c8ff45615ce5f2e"]
	slis = ["latency"]
	time_interval_type = "rolling"
	duration_in_days = 30

	rules {
		name = "breached_error_budget"
	}
	
	rules {
		name = "unhealthy_slo"
		threshold = 1
	}

	notify {
		user_ids = ["6113b0ffe4d98ae048c37010"]
	}
	
	team_id = "61443b953ffd52818bf1616a"
}
	`, sloName)
}

func testAccResourceSloConfig_update(sloName string) string {
	return fmt.Sprintf(`

resource "squadcast_slo" "test" {
	name = "%s"
	description = "Tracks some slo for some test service"
	target_slo = 99.99
	service_ids = ["6257a8eb3c8ff45615ce5f2e"]
	slis = ["latency","high-err-rate"]
	time_interval_type = "rolling"
	duration_in_days = 7

	rules {
		name = "breached_error_budget"
	}

	rules {
		name = "unhealthy_slo"
		threshold = 2
	}
	
	rules {
		name = "remaining_error_budget"
		threshold = 11
	}

	
	notify {
		user_ids = ["6113b0ffe4d98ae048c37010", "61305a78127c63c6d2c8f746"]
	}

	team_id = "61443b953ffd52818bf1616a"
}
	`, sloName)
}
