package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-test/deep"
)

type Something struct {
	First  int
	Second string
}

func TestJSONDecoder(t *testing.T) {
	const json = `{ "First": 15, "Second": "Some string" }`
	buffer := bytes.NewBuffer([]byte(json))

	expected := Something{First: 15, Second: "Some string"}

	var result Something
	err := JSONDecoder(buffer, &result)
	if err != nil {
		t.Error(err)
	}

	if diff := deep.Equal(expected, result); diff != nil {
		t.Errorf("Wrong result: %v", diff)
	}
}

func TestJSONDecoderBadJSON(t *testing.T) {
	const badJSON = `{ First: 15, "Second": "Some string" }`
	buffer := bytes.NewBuffer([]byte(badJSON))

	expected := Something{First: 15, Second: "Some string"}

	var result Something
	err := JSONDecoder(buffer, &result)
	if err == nil {
		t.Error("There should be an error")
	}

	if diff := deep.Equal(expected, result); diff == nil {
		t.Error("Structures shouldn't match")
	}
}

func fakeServer(statusCode int, data string) *httptest.Server {
	function := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", "text")
		fmt.Fprint(w, data)
	}

	return httptest.NewServer(http.HandlerFunc(function))
}

func TestGET(t *testing.T) {
	expected := "Some text"
	server := fakeServer(http.StatusOK, expected)
	defer server.Close()

	data, err := GET(server.URL)
	if data == nil {
		t.Error("Data shouldn't be nil")
	}

	if err != nil {
		t.Error(err)
	}

	result := data.String()

	if expected != result {
		t.Errorf("Wrong result, %v", result)
	}
}

func TestGET206(t *testing.T) {
	expected := "Some text"
	server := fakeServer(http.StatusPartialContent, expected)
	defer server.Close()

	data, err := GET(server.URL)
	if data == nil {
		t.Error("Data shouldn't be nil")
	}

	if err != nil {
		t.Error(err)
	}

	result := data.String()

	if expected != result {
		t.Errorf("Wrong result, %v", result)
	}
}

func TestGET404(t *testing.T) {
	expected := "Some text"
	server := fakeServer(http.StatusNotFound, expected)
	defer server.Close()

	data, err := GET(server.URL)
	if data == nil {
		t.Error("Data shouldn't be nil")
	}

	if err == nil {
		t.Error("There should be an error!")
	}

	result := data.String()

	if expected != result {
		t.Errorf("Wrong result, %v", result)
	}
}

func TestGETNoServer(t *testing.T) {
	data, err := GET("/")
	if data != nil {
		t.Error("Data should be nil!")
	}

	if err == nil {
		t.Error("There should be an error!")
	}
}
