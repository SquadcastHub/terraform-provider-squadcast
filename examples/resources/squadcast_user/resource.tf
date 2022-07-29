resource "squadcast_user" "test" {
  first_name = "test"
  last_name  = "lastname"
  email      = "test@example.com"
  role       = "stakeholder"
  abilities  = ["manage-billing"]
}