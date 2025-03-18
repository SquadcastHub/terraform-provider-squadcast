package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceWebform(t *testing.T) {
	serviceName := "webform"

	resourceName := "data.squadcast_webform.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWebformDataSourceConfigWithInputFields(serviceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "61305a9e127c63c6d2c8f76d"),
					resource.TestCheckResourceAttr(resourceName, "name", serviceName),
					resource.TestCheckResourceAttr(resourceName, "owner.0.id", "6112b0ffe4d98ae048c37010"),
					resource.TestCheckResourceAttr(resourceName, "owner.0.type", "squad"),
					resource.TestCheckResourceAttr(resourceName, "owner.0.name", "Sample Squad"),
					resource.TestCheckResourceAttr(resourceName, "header", "test header"),
					resource.TestCheckResourceAttr(resourceName, "title", "test title"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "footer_text", "test footer"),
					resource.TestCheckResourceAttr(resourceName, "footer_link", "https://www.squadcast.com"),
					resource.TestCheckResourceAttr(resourceName, "input_field.0.label", "severity"),
					resource.TestCheckResourceAttr(resourceName, "input_field.0.options.0", "critical"),
					resource.TestCheckResourceAttr(resourceName, "services.0.service_id", "6389ba2ec31b7df1caecd579"),
					resource.TestCheckResourceAttr(resourceName, "services.0.name", "Test"),
					resource.TestCheckResourceAttr(resourceName, "email_on.0", "triggered"),
				),
			},
			{
				Config: testAccWebformDataSourceConfigWithSeverity(serviceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "61305a9e127c63c6d2c8f76d"),
					resource.TestCheckResourceAttr(resourceName, "name", serviceName),
					resource.TestCheckResourceAttr(resourceName, "owner.0.id", "6112b0ffe4d98ae048c37010"),
					resource.TestCheckResourceAttr(resourceName, "owner.0.type", "squad"),
					resource.TestCheckResourceAttr(resourceName, "owner.0.name", "Sample Squad"),
					resource.TestCheckResourceAttr(resourceName, "header", "test header"),
					resource.TestCheckResourceAttr(resourceName, "title", "test title"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "footer_text", "test footer"),
					resource.TestCheckResourceAttr(resourceName, "footer_link", "https://www.squadcast.com"),
					resource.TestCheckResourceAttr(resourceName, "services.0.service_id", "6389ba2ec31b7df1caecd579"),
					resource.TestCheckResourceAttr(resourceName, "services.0.name", "Test"),
					resource.TestCheckResourceAttr(resourceName, "email_on.0", "triggered"),
				),
			},
		},
	})
}

func testAccWebformDataSourceConfigWithInputFields(serviceName string) string {
	return fmt.Sprintf(`
		resource "squadcast_webform" "test_parent" {
			name = "%s"
			team_id = "61305a9e127c63c6d2c8f76d"
			owner {
				id = "6112b0ffe4d98ae048c37010"
				type = "squad"
				name = "Sample Squad"
			}
			header = "test header"
			title = "test title"
			description = "test description"
			footer_text = "test footer"
			footer_link = "https://www.squadcast.com"
			input_field {
				label = "severity"
				options = ["critical"]
			}
			services {
				service_id = "6389ba2ec31b7df1caecd579"
				name = "Test"
			}
			email_on = ["triggered"]
		}

		data "squadcast_webform" "test" {
			name = "%s"
			team_id = "61305a9e127c63c6d2c8f76d"
		}
	`, serviceName, serviceName)
}

func testAccWebformDataSourceConfigWithSeverity(serviceName string) string {
	return fmt.Sprintf(`
		resource "squadcast_webform" "test_parent" {
			name = "%s"
			team_id = "61305a9e127c63c6d2c8f76d"
			owner {
				id = "6112b0ffe4d98ae048c37010"
				type = "squad"
				name = "Sample Squad"
			}
			header = "test header"
			title = "test title"
			description = "test description"
			footer_text = "test footer"
			footer_link = "https://www.squadcast.com"
			severity {
				type = "critical"
				description = "test description"
			}
			services {
				service_id = "6389ba2ec31b7df1caecd579"
				name = "Test"
			}
			email_on = ["triggered"]
		}

		data "squadcast_webform" "test" {
			name = "%s"
			team_id = "61305a9e127c63c6d2c8f76d"
		}
	`, serviceName, serviceName)
}
