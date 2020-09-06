package types

// AccessToken api response struct
type AccessToken struct {
	Type         string `json:"type"`
	AccessToken  string `json:"access_token"`
	IssuedAt     int64  `json:"issued_at"`
	ExpiresAt    int64  `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
}

// Meta holds the status of the request informations
type Meta struct {
	Meta struct {
		Status  int    `json:"status_code"`
		Message string `json:"error_message,omitempty"`
	} `json:"meta,omitempty"`
}
