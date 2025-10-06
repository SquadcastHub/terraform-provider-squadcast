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
	"time"

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
	Status       int           `json:"status"`
	Message      string        `json:"error_message,omitempty"`
	ConflictData interface{}   `json:"conflict_data,omitempty"`
	ErrorDetails *ErrorDetails `json:"error_details,omitempty"`
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
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

func buildErrorMessage(meta AppError, method, url string) error {
	if meta.ConflictData == nil {
		return fmt.Errorf("%s %s returned an error:\n%s", method, url, meta.Error())
	}

	var conflictItems []string

	switch data := meta.ConflictData.(type) {
	case map[string]interface{}:
		for key, value := range data {
			v := reflect.ValueOf(value)
			switch v.Kind() {
			case reflect.Slice, reflect.Map:
				if v.Len() > 0 {
					conflictItems = append(conflictItems, fmt.Sprintf("%s: %v", toHumanReadable(key), value))
				}
			case reflect.Int:
				if v.Int() > 0 {
					conflictItems = append(conflictItems, fmt.Sprintf("%s: %v", toHumanReadable(key), value))
				}
			}
		}
	case []interface{}:
		for _, item := range data {
			if itemMap, ok := item.(map[string]interface{}); ok {
				name, _ := itemMap["name"].(string)
				id, _ := itemMap["id"].(string)

				switch {
				case name != "" && id != "":
					conflictItems = append(conflictItems, fmt.Sprintf("%s (%s)", name, id))
				case name != "":
					conflictItems = append(conflictItems, name)
				case id != "":
					conflictItems = append(conflictItems, id)
				}
			}
		}
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s %s returned an error:\n%s", method, url, meta.Error()))
	if len(conflictItems) > 0 {
		builder.WriteString("\n\nConflict Details:")
		for _, item := range conflictItems {
			builder.WriteString(fmt.Sprintf("\n  • %s", item))
		}
	}

	return errors.New(builder.String())
}

func Request[TReq any, TRes any](method string, url string, client *Client, ctx context.Context, payload *TReq) (*TRes, error) {
	const maxAttempts = 3

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Build request each attempt to avoid reusing body readers
		var req *http.Request
		var err error

		if method == http.MethodGet {
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
			// Network or transport error - retry
			lastErr = err
			if attempt == maxAttempts {
				break
			}
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second // 1s,2s
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
				continue
			}
		}

		var response struct {
			Data *TRes `json:"data"`
			*Meta
		}

		func() {
			defer func() {
				if resp != nil && resp.Body != nil {
					resp.Body.Close()
				}
			}()

			bytes, readErr := io.ReadAll(resp.Body)
			if readErr != nil {
				lastErr = readErr
				return
			}

			if len(bytes) == 0 {
				if resp.StatusCode > 299 {
					lastErr = fmt.Errorf("%s %s returned an unexpected error with no body", method, url)
					return
				}
				// Success with empty body
				response.Data = nil
				lastErr = nil
				// do not retry
				// return from closure
				return
			}

			if unmarshalErr := json.Unmarshal(bytes, &response); unmarshalErr != nil {
				lastErr = unmarshalErr
				return
			}

			if resp.StatusCode > 299 {
				// Decide if we retry based on status code
				shouldRetry := resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusTooManyRequests ||
					resp.StatusCode == http.StatusRequestTimeout ||
					(resp.StatusCode >= 500 && resp.StatusCode != http.StatusNotImplemented) // 5xx except 501

				if response.Meta != nil {
					lastErr = buildErrorMessage(response.Meta.Meta, method, url)
				} else {
					lastErr = fmt.Errorf("%s %s returned %d with an unexpected error", method, url, resp.StatusCode)
				}

				if shouldRetry && attempt < maxAttempts {
					backoff := time.Duration(1<<uint(attempt-1)) * time.Second
					select {
					case <-ctx.Done():
						return
					case <-time.After(backoff):
						// proceed to next attempt
						return
					}
				}
				return
			}

			lastErr = nil
		}()

		if lastErr == nil {
			return response.Data, nil
		}

		// If we reach here and attempts remain
		// loop will retry unless context cancelled
		if attempt == maxAttempts {
			break
		}
	}

	return nil, lastErr
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
