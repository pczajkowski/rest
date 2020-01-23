package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Something struct {
	First  int
	Second string
}

func TestJSONDecoder(t *testing.T) {
	const json = `{ "First": 15, "Second": "Some string" }`
	reader := strings.NewReader(json)
	readcloser := ioutil.NopCloser(reader)
	defer readcloser.Close()

	expected := Something{First: 15, Second: "Some string"}

	var result Something
	err := JSONDecoder(readcloser, &result)
	if err != nil {
		t.Error(err)
	}

	if expected.First != result.First || expected.Second != result.Second {
		t.Errorf("Wrong result: %v", result)
	}
}

func TestJSONDecoderBadJSON(t *testing.T) {
	const badJSON = `{ First: 15, "Second": "Some string" }`
	reader := strings.NewReader(badJSON)
	readcloser := ioutil.NopCloser(reader)
	defer readcloser.Close()

	expected := Something{First: 15, Second: "Some string"}

	var result Something
	err := JSONDecoder(readcloser, &result)
	if err == nil {
		t.Error("There should be an error")
	}

	if expected.First == result.First || expected.Second == result.Second {
		t.Errorf("There should be an error on decoding, %v", result)
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
	defer data.Close()
	if err != nil {
		t.Error(err)
	}

	result, errReading := ioutil.ReadAll(data)
	if errReading != nil {
		t.Error(err)
	}

	resultString := string(result)

	if expected != resultString {
		t.Errorf("Wrong result, %v", resultString)
	}
}
