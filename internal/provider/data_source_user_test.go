package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/squadcast/terraform-provider-squadcast/internal/testdata"
)

func TestAccDataSourceUser(t *testing.T) {
	user := testdata.RandomUser()

	resourceName := "data.squadcast_user.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig(user),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", "Dheeraj"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Kumar"),
					resource.TestCheckResourceAttr(resourceName, "name", "Dheeraj Kumar"),
					resource.TestCheckResourceAttr(resourceName, "email", "dheeraj@squadcast.com"),
					resource.TestCheckResourceAttr(resourceName, "is_email_verified", "true"),
					resource.TestCheckResourceAttr(resourceName, "phone", ""),
					resource.TestCheckResourceAttr(resourceName, "is_phone_verified", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_override_dnd_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "role", "account_owner"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "Asia/Calcutta"),
					resource.TestCheckResourceAttr(resourceName, "abilities.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "abilities.0", "manage-billing"),
					resource.TestCheckResourceAttr(resourceName, "abilities.1", "manage-api-tokens"),
					resource.TestCheckResourceAttr(resourceName, "abilities.2", "manage-extensions"),
					resource.TestCheckResourceAttr(resourceName, "notification_rules.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "notification_rules.0.type", "Email"),
					resource.TestCheckResourceAttr(resourceName, "notification_rules.0.delay_minutes", "0"),
					resource.TestCheckResourceAttr(resourceName, "notification_rules.1.type", "Push"),
					resource.TestCheckResourceAttr(resourceName, "notification_rules.1.delay_minutes", "0"),
					resource.TestCheckResourceAttr(resourceName, "notification_rules.2.type", "SMS"),
					resource.TestCheckResourceAttr(resourceName, "notification_rules.2.delay_minutes", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification_rules.3.type", "Phone"),
					resource.TestCheckResourceAttr(resourceName, "notification_rules.3.delay_minutes", "2"),
					resource.TestCheckResourceAttr(resourceName, "oncall_reminder_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "oncall_reminder_rules.0.type", "Email"),
					resource.TestCheckResourceAttr(resourceName, "oncall_reminder_rules.0.delay_minutes", "60"),
				),
			},
		},
	})
}

func testAccUserDataSourceConfig(user testdata.User) string {
	return fmt.Sprintf(`
resource "squadcast_user" "test" {
	first_name = "%s"
	last_name = "%s"
	email = "%s"
	role = "stakeholder"

	abilities = ["manage-billing"]
}

data "squadcast_user" "test" {
	email = squadcast_user.test.email
}
	`, user.FirstName, user.LastName, user.Email)
}
