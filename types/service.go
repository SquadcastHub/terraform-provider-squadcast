package types

// create service api response struct
type ServiceRes struct {
	Data struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
}