data "squadcast_team" "example_team" {
  name = "example test name"
}

data "squadcast_user" "example_user" {
  email = "test@example.com"
}

data "squadcast_service" "example_service" {
  name    = "example service name"
  team_id = data.squadcast_team.example_team.id
}

data "squadcast_service" "example_service_2" {
  name    = "example service name 2"
  team_id = data.squadcast_team.example_team.id
}

resource "squadcast_webform" "example_webform" {
  name    = "example webform name"
  team_id = data.squadcast_team.example_team.id
  owner {
    type = "user"
    id   = data.squadcast_user.example_user.id
  }
  services {
    service_id = data.squadcast_service.example_service.id
    alias      = "example service alias"
  }
  services {
    service_id = data.squadcast_service.example_service_2.id
  }
  custom_domain_name = "example.com"
  header             = "formHeader"
  description        = "formDescription"
  title              = "formTitle"
  footer_text        = "footerText"
  footer_link        = "footerLink"
  email_on           = ["acknowledged", "resolved", "triggered"]
  input_field {
    label = "test_label"
    options = [
      "test1",
      "test2",
    ]
  }
  input_field {
    label = "test_label2"
    options = [
      "test1",
    ]
  }
  tags = {
    tagKey  = "tagValue"
    tagKey2 = "tagValue2"
  }
}
