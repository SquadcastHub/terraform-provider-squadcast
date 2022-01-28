package squadcast

import "github.com/terraform-provider-squadcast/types"

// Config has squadcast provider configuration
type Config struct {
	AccessToken string
	// IssuedAt     int64
	// ExpiresAt    int64
	// RefreshToken string
}

var (
	squadcastAPIHost    string = "https://api.squadcast.com"
	squadcastAPIVersion string = "/v3"
	accessToken         *types.AccessToken
)

var apiEndpoints = map[string]string{"US": "https://api.squadcast.com", "EU": "https://api.eu.squadcast.com"}
