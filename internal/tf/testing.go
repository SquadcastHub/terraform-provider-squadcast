package tf

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func stateAttr(s *terraform.State, resourceName, key string) (string, bool) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourceName {
			continue
		}

		val, ok := rs.Primary.Attributes[key]
		if !ok {
			return "", false
		}

		return val, true
	}

	return "", false
}

func StateAttr(s *terraform.State, resourceName, key string) (string, error) {
	val, found := stateAttr(s, resourceName, key)
	if !found {
		return "", fmt.Errorf("cannot find %s.%s in the state", resourceName, key)
	}

	return val, nil
}
