package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getDetailedError(data []byte, err error) error {
	var limit int64 = 20

	switch detailed := err.(type) {
	case *json.SyntaxError:
		start := detailed.Offset - 1
		if start < 0 {
			start = 0
		}

		end := start + limit
		dataLength := int64(len(data))
		if end > dataLength {
			end = dataLength
		}

		badPart := string(data[start:end])
		return fmt.Errorf("%s:\n%s", err.Error(), badPart)

	case *json.UnmarshalTypeError:
		start := detailed.Offset - limit
		if start < 0 {
			start = 0
		}

		badPart := string(data[start:detailed.Offset])
		return fmt.Errorf("%s:\n%s", err.Error(), badPart)
	}

	return err
}

//JSONDecoder decodes json from given bytes array to target object.
func JSONDecoder(data []byte, target interface{}) error {
	err := json.Unmarshal(data, target)
	if err != nil {
		err = getDetailedError(data, err)
	}

	return err
}

//GET returns bytes buffer of response body.
func GET(url string) (*bytes.Buffer, error) {
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, fmt.Errorf("Response error: %v", err)
	}

	body, err := BodyToBuffer(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading body: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return body, fmt.Errorf("Request unsuccessful: %v - %v", resp.Status, url)
	}

	return body, nil
}

//HEAD returns map[string][]string of response header.
func HEAD(url string) (map[string][]string, error) {
	resp, err := http.Head(url)
	if err != nil {
		return nil, fmt.Errorf("Response error: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return resp.Header, fmt.Errorf("Request unsuccessful: %v - %v", resp.Status, url)
	}

	return resp.Header, nil
}

//BodyToBuffer reads data from ReadCloser and returns bytes buffer.
func BodyToBuffer(data io.ReadCloser) (*bytes.Buffer, error) {
	var buffer bytes.Buffer

	_, err := buffer.ReadFrom(data)
	if err != nil {
		return nil, err
	}

	return &buffer, nil
}
