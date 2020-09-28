package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// PostRequest performs a HTTP Post Request
func PostRequest(url string, jsonBody interface{}, response interface{}) error {
	serializedBody, err := json.Marshal(jsonBody)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(serializedBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, response)
}
