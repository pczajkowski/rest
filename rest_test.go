package rest

import (
	"io/ioutil"
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
