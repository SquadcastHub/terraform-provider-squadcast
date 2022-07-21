package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	Host   string
	Region string

	RefreshToken   string
	AccessToken    string
	OrganizationID string

	UserAgent        string
	BaseURLV2        string
	BaseURLV3        string
	AuthBaseURL      string
	IngestionBaseURL string
}

type ErrorDetails struct {
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
	Link        string `json:"link,omitempty"`
	Errors      any    `json:"errors,omitempty"`
}

type AppError struct {
	Status       int           `json:"status"`
	Message      string        `json:"error_message,omitempty"`
	ConflictData *any          `json:"conflict_data,omitempty"`
	ErrorDetails *ErrorDetails `json:"error_details,omitempty"`
}

func (err *AppError) Error() string {
	str := fmt.Sprintf("[%d] %s", err.Status, err.Message)
	if err.ErrorDetails != nil {
		str += fmt.Sprintf("\ndetails: %#v", err.ErrorDetails)
	}
	return str
}

// Meta holds the status of the request informations
type Meta struct {
	Meta AppError `json:"meta,omitempty"`
}

func Request[TReq any, TRes any](method string, url string, client *Client, ctx context.Context, payload *TReq) (*TRes, error) {
	var req *http.Request
	var err error

	if method == "GET" {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	} else {
		buf := &bytes.Buffer{}
		if payload != nil {
			body, err := json.Marshal(payload)
			if err != nil {
				return nil, err
			}
			buf = bytes.NewBuffer(body)
		}
		req, err = http.NewRequestWithContext(ctx, method, url, buf)
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.AccessToken))
	req.Header.Set("User-Agent", client.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data *TRes `json:"data"`
		*Meta
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		if resp.StatusCode > 299 {
			return nil, fmt.Errorf("%s %s returned an unexpected error with no body", method, url)
		} else {
			return nil, nil
		}
	}

	if err := json.Unmarshal(bytes, &response); err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		if response.Meta != nil {
			return nil, fmt.Errorf("%s %s returned an error:\n%s", method, url, response.Meta.Meta.Error())
		} else {
			return nil, fmt.Errorf("%s %s returned %d with an unexpected error: %#v", method, url, resp.StatusCode, response)
		}
	}

	return response.Data, nil
}

func RequestSlice[TReq any, TRes any](method string, url string, client *Client, ctx context.Context, payload *TReq) ([]*TRes, error) {
	data, err := Request[TReq, []*TRes](method, url, client, ctx, payload)
	if err != nil {
		return nil, err
	}

	return *data, nil
}

func IsResourceNotFoundError(e error) bool {
	return strings.Contains(e.Error(), "[404]")
}
