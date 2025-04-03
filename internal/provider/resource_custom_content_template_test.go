package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
)

func TestAccResourceCustomContentTemplate(t *testing.T) {
	resourceName := "squadcast_custom_content_template.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCustomContentTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCustomContentTemplateConfig_defaults(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "service_id", "61361611c2fc70c3101ca7dd"),
					resource.TestCheckResourceAttr(resourceName, "alert_source", "APImetrics"),
					resource.TestCheckResourceAttr(resourceName, "message_template", "{{.labels.alertname}}-{{.labels.deployment}}"),
					resource.TestCheckResourceAttr(resourceName, "description_template", "{{.labels.description}}"),
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

func testAccCheckCustomContentTemplateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_custom_content_template" {
			continue
		}

		_, err := client.GetCustomContentTemplateOverlay(context.Background(), rs.Primary.Attributes["service_id"], rs.Primary.Attributes["alert_source_shortname"])
		if err != nil {
			return err
		}
		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceCustomContentTemplateConfig_defaults() string {
	return fmt.Sprintf(`
resource "squadcast_custom_content_template" "test" {
	service_id = "61361611c2fc70c3101ca7dd"
	message_template     = "{{.labels.alertname}}-{{.labels.deployment}}"
	description_template = "{{.labels.description}}"
	alert_source         = "APImetrics"
}
	`)
}
