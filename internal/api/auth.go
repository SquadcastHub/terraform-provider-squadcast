package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type AccessToken struct {
	Type         string `json:"type"`
	AccessToken  string `json:"access_token"`
	IssuedAt     int64  `json:"issued_at"`
	ExpiresAt    int64  `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
}

func (client *Client) GetAccessToken(ctx context.Context) (*AccessToken, error) {
	path := "/oauth/access-token"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.AuthBaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Refresh-Token", client.RefreshToken)
	req.Header.Set("User-Agent", client.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data AccessToken `json:"data"`
		*Meta
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &response); err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, errors.New(response.Meta.Meta.Message)
	}

	return &response.Data, nil
}
