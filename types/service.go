package types

// ServiceRes is for unmarshalling create service api response
type ServiceRes struct {
	Data struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
}
