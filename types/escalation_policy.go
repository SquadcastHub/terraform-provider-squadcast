package types

// create escalation-policy api response struct
type EscalationPolicyRes struct {
	Data []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
}
