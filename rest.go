package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// JSONDecoder decodes json from given io.ReadCloser to target object.
func JSONDecoder(data io.ReadCloser, target interface{}) error {
	decoder := json.NewDecoder(data)

	err := decoder.Decode(target)
	if err != nil {
		return err
	}

	return nil
}

// GET returns io.ReadCloser of response body. Don't forget to close it.
func GET(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Response error: %v", err)
	}

	if resp.StatusCode != 200 {
		return resp.Body, fmt.Errorf("Request unsuccessful: %v (%v) - %v", resp.Status, resp.StatusCode, url)
	}

	return resp.Body, nil
}
