package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// JSONDecoder decodes json from given io.ReadCloser to target object.
func JSONDecoder(data io.ReadCloser, target interface{}) error {
	decoder := json.NewDecoder(data)

	err := decoder.Decode(target)

	return err
}

// GET returns io.ReadCloser of response body. Don't forget to close it.
func GET(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Response error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return resp.Body, fmt.Errorf("Request unsuccessful: %v - %v", resp.Status, url)
	}

	return resp.Body, nil
}

// BodyToString returns string read from given Body.
func BodyToString(data io.ReadCloser) (string, error) {
	var buffer bytes.Buffer

	_, err := buffer.ReadFrom(data)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
