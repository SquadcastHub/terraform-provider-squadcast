package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRunbook(t *testing.T) {
	runbookName := acctest.RandomWithPrefix("runbook")

	resourceName := "data.squadcast_runbook.test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRunbookDataSourceConfig(runbookName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "team_id", "613611c1eb22db455cfa789f"),
					resource.TestCheckResourceAttr(resourceName, "name", runbookName),
					resource.TestCheckResourceAttr(resourceName, "entity_owner.id", "6113b0ffe4d98ae048c37010"),
					resource.TestCheckResourceAttr(resourceName, "entity_owner.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "steps.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "steps.0.content", "some text here"),
					resource.TestCheckResourceAttr(resourceName, "steps.1.content", "some text here 2"),
				),
			},
		},
	})
}

func testAccRunbookDataSourceConfig(runbookName string) string {
	return fmt.Sprintf(`
resource "squadcast_runbook" "test" {
	name = "%s"
	team_id = "613611c1eb22db455cfa789f"
	
	entity_owner{
		id = "6113b0ffe4d98ae048c37010"
		type = "user"
	}

	steps {
		content = "some text here"
	}

	steps {
		content = "some text here 2"
	}
}

data "squadcast_runbook" "test" {
	name = squadcast_runbook.test.name
	team_id = "613611c1eb22db455cfa789f"
}
	`, runbookName)
}
