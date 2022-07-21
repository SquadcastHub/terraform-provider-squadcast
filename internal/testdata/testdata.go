package testdata

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

type User struct {
	Email     string
	FirstName string
	LastName  string
}

func RandomUser() User {
	firstName := fmt.Sprintf("testuser%s", acctest.RandStringFromCharSet(10, "abcdefghijlkmnopqrstuvwxyz"))
	email := firstName + "@example.com"

	return User{
		Email:     email,
		FirstName: firstName,
		LastName:  "lastname",
	}
}
