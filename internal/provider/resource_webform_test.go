package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

func TestAccResourceWebform(t *testing.T) {
	webformName := "webform"
	resourceName := "squadcast_webform.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckWebformDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWebformConfigWithInputFields(webformName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "61305a9e127c63c6d2c8f76d"),
					resource.TestCheckResourceAttr(resourceName, "name", webformName),
					resource.TestCheckResourceAttr(resourceName, "owner.0.id", "6112b0ffe4d98ae048c37010"),
					resource.TestCheckResourceAttr(resourceName, "owner.0.type", "squad"),
					resource.TestCheckResourceAttr(resourceName, "owner.0.name", "Sample Squad"),
					resource.TestCheckResourceAttr(resourceName, "header", "test header"),
					resource.TestCheckResourceAttr(resourceName, "title", "test title"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "footer_text", "test footer"),
					resource.TestCheckResourceAttr(resourceName, "footer_link", "https://www.squadcast.com"),
					resource.TestCheckResourceAttr(resourceName, "input_field.0.label", "severity"),
					resource.TestCheckResourceAttr(resourceName, "input_field.0.options.0", "critical"),
					resource.TestCheckResourceAttr(resourceName, "services.0.service_id", "6389ba2ec31b7df1caecd579"),
					resource.TestCheckResourceAttr(resourceName, "services.0.name", "Test"),
				),
			},
			{
				Config: testAccResourceWebformConfigWithSeverity(webformName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "61305a9e127c63c6d2c8f76d"),
					resource.TestCheckResourceAttr(resourceName, "name", webformName),
					resource.TestCheckResourceAttr(resourceName, "owner.0.id", "6112b0ffe4d98ae048c37010"),
					resource.TestCheckResourceAttr(resourceName, "owner.0.type", "squad"),
					resource.TestCheckResourceAttr(resourceName, "owner.0.name", "Sample Squad"),
					resource.TestCheckResourceAttr(resourceName, "header", "test header"),
					resource.TestCheckResourceAttr(resourceName, "title", "test title"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "footer_text", "test footer"),
					resource.TestCheckResourceAttr(resourceName, "footer_link", "https://www.squadcast.com"),
					resource.TestCheckResourceAttr(resourceName, "severity.0.type", "critical"),
					resource.TestCheckResourceAttr(resourceName, "severity.0.description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "services.0.service_id", "6389ba2ec31b7df1caecd579"),
					resource.TestCheckResourceAttr(resourceName, "services.0.name", "Test"),
				),
			},
			{
				Config: testAccResourceWebformConfig_update(webformName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "61305a9e127c63c6d2c8f76d"),
					resource.TestCheckResourceAttr(resourceName, "name", webformName),
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
				Config: testAccResourceWebformConfig_tags(webformName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "61305a9e127c63c6d2c8f76d"),
					resource.TestCheckResourceAttr(resourceName, "name", webformName),
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
					resource.TestCheckResourceAttr(resourceName, "tags.testKey", "testVal"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					teamID, err := tf.StateAttr(s, "squadcast_webform", "team_id")
					if err != nil {
						return "", err
					}

					return teamID + ":" + webformName, nil
				},
			},
		},
	})
}

func testAccCheckWebformDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_webform" {
			continue
		}

		_, err := client.GetWebformById(context.Background(), rs.Primary.Attributes["team_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("expected webform to be destroyed, %s found", rs.Primary.ID)
		}

		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceWebformConfigWithInputFields(webformName string) string {
	return fmt.Sprintf(`
		resource "squadcast_webform" "test" {
			name = "%s"
			team_id = "61305a9e127c63c6d2c8f76d"
			owner {
				id = "6112b0ffe4d98ae048c37010"
				type = "squad"
				name = "Sample Squad"
			}
			header = "test header"
			title = "test title"
			description = ""
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
		}
	`, webformName)
}

func testAccResourceWebformConfigWithSeverity(webformName string) string {
	return fmt.Sprintf(`
		resource "squadcast_webform" "test" {
			name = "%s"
			team_id = "61305a9e127c63c6d2c8f76d"
			owner {
				id = "6112b0ffe4d98ae048c37010"
				type = "squad"
				name = "Sample Squad"
			}
			header = "test header"
			title = "test title"
			description = ""
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
		}
	`, webformName)
}

func testAccResourceWebformConfig_update(webformName string) string {
	return fmt.Sprintf(`
		resource "squadcast_webform" "test" {
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
	`, webformName)
}

func testAccResourceWebformConfig_tags(webformName string) string {
	return fmt.Sprintf(`
		resource "squadcast_webform" "test" {
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
			tags = {
				"testKey" = "testVal"
			}
		}
	`, webformName)
}
