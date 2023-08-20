data "squadcast_team" "team" {
  name = "Default Team"
}

resource "squadcast_status_page" "test_status_page" {
	team_id = data.squadcast_team.team.id
	name = "Test Status Page"
	description = "Status Page description"
	is_public = true
	domain_name = "test-statuspage"
	timezone = "Asia/Kolkata"
	contact_email = "example@test.com"
	theme_color {
		primary = "#000000"
		secondary = "#dddddd"
	}
	owner {
		type = "team"
		id = data.squadcast_team.team.id
	}
}

resource "squadcast_status_page_component" "example_component" {
	status_page_id = squadcast_status_page.test_status_page.id
	name = "Component 1"
	description = "Component 1 description"
}

resource "squadcast_status_page_component" "example_component_2" {
	status_page_id = squadcast_status_page.test_status_page.id
	name = "Component 2"
	description = "Component 2 description"
}

resource "squadcast_status_page_group" "example_group" {
  	status_page_id = squadcast_status_page.test_status_page.id
	name = "Group 1"
	description = "Group 1 description"
	allow_subscription = true
	component_ids = [ 
		squadcast_status_page_component.example_component.id,
        squadcast_status_page_component.example_component_2.id
	 ]
}
