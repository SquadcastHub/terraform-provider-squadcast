package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/terraform-provider-squadcast/types"
)

type APICallError struct {
	StatusCode int
	Headers    http.Header

	err error
}

func (a APICallError) Error() string {
	return a.err.Error()
}

// Request struct contains all data for a simple http
// request
type Request struct {
	method  string
	url     string
	json    interface{}
	headers map[string]string
	to      interface{}
	failure interface{}
}

func New() *Request {
	return &Request{
		method: http.MethodGet,
		url:    "",
		json:   types.JSON{},
		headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

// Method sets the method and url parameters of the request.
func (r *Request) Method(method, url string) *Request {
	r.method = method
	r.url = url

	return r
}

// Get configures the request to GET the passed URL.
func (r *Request) Get(url string) *Request {
	return r.Method(http.MethodGet, url)
}

// Head configures the request to HEAD (HEAD request...) the passed URL.
func (r *Request) Head(url string) *Request {
	return r.Method(http.MethodHead, url)
}

// Post configures the request to POST to the passed URL.
func (r *Request) Post(url string) *Request {
	return r.Method(http.MethodPost, url)
}

// Put configures the request to PUT to the passed URL.
func (r *Request) Put(url string) *Request {
	return r.Method(http.MethodPut, url)
}

// Patch configures the request to PATCH to the passed URL.
func (r *Request) Patch(url string) *Request {
	return r.Method(http.MethodPatch, url)
}

// Delete configures the request to DELETE to the passed URL.
func (r *Request) Delete(url string) *Request {
	return r.Method(http.MethodDelete, url)
}

// Options configures the request to OPTIONS to the passed URL.
func (r *Request) Options(url string) *Request {
	return r.Method(http.MethodOptions, url)
}

// Data sets the body data to be sent in the request
func (r *Request) Data(body interface{}) *Request {
	r.json = body
	return r
}

// SetHeader is used to set key value pair in the request header
func (r *Request) SetHeader(key, value string) *Request {
	r.headers[key] = value
	return r
}

// With records the passed variable (pointer) as a destination for any
// response decoding
func (r *Request) With(decoder interface{}) *Request {
	r.to = decoder
	return r
}

// WithFail records the passed variable (pointer) as a destination for any
// response decoding when the API request fails (based on the status code)
func (r *Request) WithFail(decoder interface{}) *Request {
	r.failure = decoder
	return r
}

// Do calls the endpoint and decodes the returned response into the
func (r *Request) Do() error {
	b := bytes.NewBuffer(nil)
	err := json.NewEncoder(b).Encode(r.json)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(r.method, r.url, b)
	if err != nil {
		return err
	}

	for key, value := range r.headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if resp.StatusCode > 299 {
		errString := fmt.Errorf("Status code error: %d", resp.StatusCode)
		if r.failure != nil {
			err := dec.Decode(r.failure)
			if err != nil {
				errString = fmt.Errorf("%s, decode error: %w", errString, err)
			}
		}
		return &APICallError{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			err:        errString,
		}
	}

	if r.to == nil {
		return nil
	}
	return dec.Decode(r.to)
}
