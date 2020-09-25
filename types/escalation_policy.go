package types

// EscalationPolicyRes is to unmarshal api response
type EscalationPolicyRes struct {
	Data []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"data"`
}
