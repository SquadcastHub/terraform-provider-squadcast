package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
)

func TestAccResourceDedupKeyOverlay(t *testing.T) {
	resourceName := "squadcast_dedup_key_overlay.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckDedupKeyOverlayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDedupKeyOverlayConfig_defaults(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "service_id", "61361611c2fc70c3101ca7dd"),
					resource.TestCheckResourceAttr(resourceName, "alert_source", "APImetrics"),
					resource.TestCheckResourceAttr(resourceName, "dedup_key_overlay_template", "Alertname: sample"),
					resource.TestCheckResourceAttr(resourceName, "duration", "100"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "6731ef12d648e43996fe552c:APImetrics",
			},
		},
	})
}

func testAccCheckDedupKeyOverlayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_dedup_key_overlay" {
			continue
		}

		_, err := client.GetDedupKeyOverlay(context.Background(), rs.Primary.Attributes["service_id"], rs.Primary.Attributes["alert_source_shortname"])
		if err != nil {
			return err
		}
		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceDedupKeyOverlayConfig_defaults() string {
	return fmt.Sprintf(`
resource "squadcast_dedup_key_overlay" "test" {
	service_id = "61361611c2fc70c3101ca7dd"
	dedup_key_overlay_template = "Alertname: sample"
	duration                   = 100
	alert_source               = "APImetrics"
}
	`)
}
