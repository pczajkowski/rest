package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// JSONDecoder decodes json from given bytes buffer to target object.
func JSONDecoder(data bytes.Buffer, target interface{}) error {
	err := json.Unmarshal(data.Bytes(), target)

	return err
}

// GET returns bytes buffer of response body.
func GET(url string) (*bytes.Buffer, error) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Response error: %v", err)
	}

	body, err := bodyToBuffer(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("Request unsuccessful: %v - %v", resp.Status, url)
	}

	return body, nil
}

func bodyToBuffer(data io.ReadCloser) (*bytes.Buffer, error) {
	var buffer bytes.Buffer

	_, err := buffer.ReadFrom(data)
	if err != nil {
		return nil, err
	}

	return &buffer, nil
}
