package squadcast

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/terraform-provider-squadcast/types"
)

var accessToken *types.AccessToken

const (
	squadcastAPIHost    string = "https://api.squadcast.com"
	squadcastAPIVersion string = "/v3"
)

func getAPIFullURL(path string) string {
	return squadcastAPIHost + squadcastAPIVersion + path
}

// getAccessToken fetches bearer access token using refresh token
func getAccessToken(refreshToken string) (string, error) {
	if accessToken == nil || accessToken.ExpiresAt <= time.Now().Unix() {
		if err := updateAccessToken(refreshToken); err != nil {
			return "", err
		}
	}

	return accessToken.Type + " " + accessToken.AccessToken, nil
}

func updateAccessToken(refreshToken string) error {
	if refreshToken == "" {
		return errors.New("refresh token required")
	}

	path := "/oauth/access-token"
	req, err := http.NewRequest(http.MethodGet, getAPIFullURL(path), nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-Refresh-Token", refreshToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response struct {
		Data types.AccessToken `json:"data"`
		*types.Meta
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		errMsg := response.Meta.Meta.Message
		if strings.TrimSpace(errMsg) == "" {
			errMsg = "Something went wrong"
		}

		return errors.New(errMsg)
	}

	accessToken = &response.Data
	return nil
}
