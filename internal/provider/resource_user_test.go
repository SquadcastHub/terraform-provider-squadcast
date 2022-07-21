package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/squadcast/terraform-provider-squadcast/internal/api"
	"github.com/squadcast/terraform-provider-squadcast/internal/testdata"
)

func TestAccResourceUserNoAbilities(t *testing.T) {
	user := testdata.RandomUser()

	resourceName := "squadcast_user.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserConfig_user_noabilities(user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", user.FirstName),
					resource.TestCheckResourceAttr(resourceName, "last_name", user.LastName),
					resource.TestCheckResourceAttr(resourceName, "email", user.Email),
					resource.TestCheckResourceAttr(resourceName, "role", "user"),
					resource.TestCheckNoResourceAttr(resourceName, "abilities.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     user.Email,
			},
		},
	})
}

func TestAccResourceUserAbilities(t *testing.T) {
	user := testdata.RandomUser()

	resourceName := "squadcast_user.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserConfig_user_abilities(user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", user.FirstName),
					resource.TestCheckResourceAttr(resourceName, "last_name", user.LastName),
					resource.TestCheckResourceAttr(resourceName, "email", user.Email),
					resource.TestCheckResourceAttr(resourceName, "role", "user"),
					resource.TestCheckResourceAttr(resourceName, "abilities.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "abilities.0", "manage-billing"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     user.Email,
			},
		},
	})
}

func TestAccResourceUserStakeholderNoAbilities(t *testing.T) {
	user := testdata.RandomUser()

	resourceName := "squadcast_user.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserConfig_stakeholder_noabilities(user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", user.FirstName),
					resource.TestCheckResourceAttr(resourceName, "last_name", user.LastName),
					resource.TestCheckResourceAttr(resourceName, "email", user.Email),
					resource.TestCheckResourceAttr(resourceName, "role", "stakeholder"),
					resource.TestCheckNoResourceAttr(resourceName, "abilities.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     user.Email,
			},
		},
	})
}

func TestAccResourceUserStakeholderAbilities(t *testing.T) {
	user := testdata.RandomUser()

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceUserConfig_stakeholder_abilities(user),
				ExpectError: regexp.MustCompile("stakeholders cannot have special abilities"),
			},
		},
	})
}

func TestAccResourceUserNoAbilitiesToStakeholderNOAbilities(t *testing.T) {
	user := testdata.RandomUser()

	resourceName := "squadcast_user.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserConfig_user_noabilities(user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", user.FirstName),
					resource.TestCheckResourceAttr(resourceName, "last_name", user.LastName),
					resource.TestCheckResourceAttr(resourceName, "email", user.Email),
					resource.TestCheckResourceAttr(resourceName, "role", "user"),
					resource.TestCheckNoResourceAttr(resourceName, "abilities.#"),
				),
			},
			{
				Config: testAccResourceUserConfig_stakeholder_noabilities(user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", user.FirstName),
					resource.TestCheckResourceAttr(resourceName, "last_name", user.LastName),
					resource.TestCheckResourceAttr(resourceName, "email", user.Email),
					resource.TestCheckResourceAttr(resourceName, "role", "stakeholder"),
					resource.TestCheckNoResourceAttr(resourceName, "abilities.#"),
				),
			},
		},
	})
}

func TestAccResourceUserNoAbilitiesToStakeholderAbilities(t *testing.T) {
	user := testdata.RandomUser()

	resourceName := "squadcast_user.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserConfig_user_noabilities(user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", user.FirstName),
					resource.TestCheckResourceAttr(resourceName, "last_name", user.LastName),
					resource.TestCheckResourceAttr(resourceName, "email", user.Email),
					resource.TestCheckResourceAttr(resourceName, "role", "user"),
					resource.TestCheckNoResourceAttr(resourceName, "abilities.#"),
				),
			},
			{
				Config:      testAccResourceUserConfig_stakeholder_abilities(user),
				ExpectError: regexp.MustCompile("stakeholders cannot have special abilities"),
			},
		},
	})
}

func TestAccResourceUserAbilitiesToStakeholderNoAbilities(t *testing.T) {
	user := testdata.RandomUser()

	resourceName := "squadcast_user.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserConfig_user_abilities(user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", user.FirstName),
					resource.TestCheckResourceAttr(resourceName, "last_name", user.LastName),
					resource.TestCheckResourceAttr(resourceName, "email", user.Email),
					resource.TestCheckResourceAttr(resourceName, "role", "user"),
					resource.TestCheckResourceAttr(resourceName, "abilities.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "abilities.0", "manage-billing"),
				),
			},
			{
				Config: testAccResourceUserConfig_stakeholder_noabilities(user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", user.FirstName),
					resource.TestCheckResourceAttr(resourceName, "last_name", user.LastName),
					resource.TestCheckResourceAttr(resourceName, "email", user.Email),
					resource.TestCheckResourceAttr(resourceName, "role", "stakeholder"),
					resource.TestCheckNoResourceAttr(resourceName, "abilities.#"),
				),
			},
		},
	})
}

func TestAccResourceUserAbilitiesToStakeholderAbilities(t *testing.T) {
	user := testdata.RandomUser()

	resourceName := "squadcast_user.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserConfig_user_abilities(user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", user.FirstName),
					resource.TestCheckResourceAttr(resourceName, "last_name", user.LastName),
					resource.TestCheckResourceAttr(resourceName, "email", user.Email),
					resource.TestCheckResourceAttr(resourceName, "role", "user"),
					resource.TestCheckResourceAttr(resourceName, "abilities.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "abilities.0", "manage-billing"),
				),
			},
			{
				Config:      testAccResourceUserConfig_stakeholder_abilities(user),
				ExpectError: regexp.MustCompile("stakeholders cannot have special abilities"),
			},
		},
	})
}

func TestAccResourceUserStakeholderNoAbilitiesToUserNoAbilities(t *testing.T) {
	user := testdata.RandomUser()

	resourceName := "squadcast_user.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserConfig_stakeholder_noabilities(user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", user.FirstName),
					resource.TestCheckResourceAttr(resourceName, "last_name", user.LastName),
					resource.TestCheckResourceAttr(resourceName, "email", user.Email),
					resource.TestCheckResourceAttr(resourceName, "role", "stakeholder"),
					resource.TestCheckNoResourceAttr(resourceName, "abilities.#"),
				),
			},
			{
				Config: testAccResourceUserConfig_user_noabilities(user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", user.FirstName),
					resource.TestCheckResourceAttr(resourceName, "last_name", user.LastName),
					resource.TestCheckResourceAttr(resourceName, "email", user.Email),
					resource.TestCheckResourceAttr(resourceName, "role", "user"),
					resource.TestCheckNoResourceAttr(resourceName, "abilities.#"),
				),
			},
		},
	})
}

func TestAccResourceUserStakeholderNoAbilitiesToUserAbilities(t *testing.T) {
	user := testdata.RandomUser()

	resourceName := "squadcast_user.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserConfig_stakeholder_noabilities(user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", user.FirstName),
					resource.TestCheckResourceAttr(resourceName, "last_name", user.LastName),
					resource.TestCheckResourceAttr(resourceName, "email", user.Email),
					resource.TestCheckResourceAttr(resourceName, "role", "stakeholder"),
					resource.TestCheckNoResourceAttr(resourceName, "abilities.#"),
				),
			},
			{
				Config: testAccResourceUserConfig_user_abilities(user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "first_name", user.FirstName),
					resource.TestCheckResourceAttr(resourceName, "last_name", user.LastName),
					resource.TestCheckResourceAttr(resourceName, "email", user.Email),
					resource.TestCheckResourceAttr(resourceName, "role", "user"),
					resource.TestCheckResourceAttr(resourceName, "abilities.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "abilities.0", "manage-billing"),
				),
			},
		},
	})
}

func testAccCheckUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "squadcast_user" {
			continue
		}

		_, err := client.GetUserById(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("expected user to be destroyed, %s found", rs.Primary.ID)
		}

		if !api.IsResourceNotFoundError(err) {
			return err
		}
	}

	return nil
}

func testAccResourceUserConfig_user_noabilities(user testdata.User) string {
	return fmt.Sprintf(`
resource "squadcast_user" "test" {
	first_name = "%s"
	last_name = "%s"
	email = "%s"
	role = "user"
}
	`, user.FirstName, user.LastName, user.Email)
}

func testAccResourceUserConfig_user_abilities(user testdata.User) string {
	return fmt.Sprintf(`
resource "squadcast_user" "test" {
	first_name = "%s"
	last_name = "%s"
	email = "%s"
	role = "user"

	abilities = ["manage-billing"]
}
	`, user.FirstName, user.LastName, user.Email)
}

func testAccResourceUserConfig_stakeholder_noabilities(user testdata.User) string {
	return fmt.Sprintf(`
resource "squadcast_user" "test" {
	first_name = "%s"
	last_name = "%s"
	email = "%s"
	role = "stakeholder"
}
	`, user.FirstName, user.LastName, user.Email)
}

func testAccResourceUserConfig_stakeholder_abilities(user testdata.User) string {
	return fmt.Sprintf(`
resource "squadcast_user" "test" {
	first_name = "%s"
	last_name = "%s"
	email = "%s"
	role = "stakeholder"

	abilities = ["manage-billing"]
}
	`, user.FirstName, user.LastName, user.Email)
}
