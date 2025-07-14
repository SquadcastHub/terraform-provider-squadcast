package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/hasura/go-graphql-client"
)

type Client struct {
	Host   string
	Region string

	RefreshToken   string
	AccessToken    string
	OrganizationID string

	UserAgent        string
	BaseURLV3        string
	BaseURLV4        string
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
	Status       int                    `json:"status"`
	Message      string                 `json:"error_message,omitempty"`
	ConflictData map[string]interface{} `json:"conflict_data,omitempty"`
	ErrorDetails *ErrorDetails          `json:"error_details,omitempty"`
}

var GraphQLClient *graphql.Client

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

func toHumanReadable(key string) string {
	words := strings.Split(key, "_")
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, " ")
}

func buildErrorMessage(meta AppError, method, url string) error {
	if meta.ConflictData == nil {
		return fmt.Errorf("%s %s returned an error:\n%s", method, url, meta.Error())
	}

	nonEmptyFields := make(map[string]interface{})
	for key, value := range meta.ConflictData {
		v := reflect.ValueOf(value)
		if (v.Kind() == reflect.Slice || v.Kind() == reflect.Map) && v.Len() > 0 {
			nonEmptyFields[key] = value
		} else if v.Kind() == reflect.Int && v.Int() > 0 {
			nonEmptyFields[key] = value
		}
	}

	errorMessage := fmt.Sprintf("%s %s returned an error:\n%s", method, url, meta.Error())
	if len(nonEmptyFields) > 0 {
		conflictDataMessage := " Conflict data details:\n"
		for key, value := range nonEmptyFields {
			humanReadableKey := toHumanReadable(key)
			conflictDataMessage += fmt.Sprintf("%s: %+v\n", humanReadableKey, value)

		}
		errorMessage += conflictDataMessage
	}
	return fmt.Errorf(errorMessage)
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
			return nil, buildErrorMessage(response.Meta.Meta, method, url)
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

// GraphQLRequest is a generic function to make graphql requests
// method values can be query/mutate
func GraphQLRequest[TReq any](method string, client *Client, ctx context.Context, payload *TReq, variables map[string]interface{}) (*TReq, error) {
	switch method {
	case "query":
		if err := GraphQLClient.WithDebug(false).Query(ctx, payload, variables); err != nil {
			return nil, err
		}
	case "mutate":
		if err := GraphQLClient.WithDebug(false).Mutate(ctx, payload, variables); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid method")
	}

	return payload, nil
}
