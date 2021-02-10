package parsers

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"gotest.tools/assert"
	"testing"
)

type simple struct {
	A int `validate:"required"`
	B int
}

type layered struct {
	C simple `validate:"-"`
	D bool
	E string `validate:"required"`
}

type complex struct {
	F layered
}

type arr struct {
	G string
	H string `validate:"required"`
}

func TestJsonToStruct(t *testing.T) {
	// case 1 simple
	input := []byte("[{\"A\": 1, \"B\": 2},{\"A\": 3, \"B\": 4}]")
	simpleWant := []simple{{1, 2}, {3, 4}}
	simpleSpec := []simple{}

	err := JsonToSpec(input, &simpleSpec)
	assert.NilError(t, err)
	assert.DeepEqual(t, &simpleWant, &simpleSpec)

	// case 2 complex
	input = []byte("{\"F\": {\"D\": true, \"E\": \"nice\", \"C\": {\"A\": 1, \"B\": 2}}}")
	complexWant := &complex{layered{simple{1, 2}, true, "nice"}}
	complexSpec := &complex{}

	err = JsonToSpec(input, complexSpec)
	assert.NilError(t, err)
	assert.DeepEqual(t, complexWant, complexSpec)
}

func TestParseReqBodyToStruct(t *testing.T) {
	request, _ := http.NewRequest("POST", "http://test.jsonparsing.com", bytes.NewReader([]byte("")))
	// case 1 simple
	request.Body = ioutil.NopCloser(bytes.NewReader([]byte("{\"A\": 1, \"B\": 2}")))
	simpleWant := &simple{1, 2}
	simpleSpec := &simple{}

	err := ParseReqBodyToSpec(request, simpleSpec)
	assert.NilError(t, err)
	assert.DeepEqual(t, simpleWant, simpleSpec)

	// case 2 complex
	request.Body = ioutil.NopCloser(bytes.NewReader([]byte("{\"F\": {\"D\": true, \"E\": \"nice\", \"C\": {\"A\": 1, \"B\": 2}}}")))
	complexWant := &complex{layered{simple{1, 2}, true, "nice"}}
	complexSpec := &complex{}

	err = ParseReqBodyToSpec(request, complexSpec)
	assert.NilError(t, err)
	assert.DeepEqual(t, complexWant, complexSpec)
}

func TestParseAndValidate(t *testing.T) {
	request, _ := http.NewRequest("POST", "http://test.jsonparsing.com", bytes.NewReader([]byte("")))

	request.Body = ioutil.NopCloser(bytes.NewReader([]byte("{\"F\": {\"D\": true, \"E\": \"\", \"C\": {\"A\": 1, \"B\": 2}}}")))
	complexSpec := &complex{}
	want := "Key: 'complex.F.E' Error:Field validation for 'E' failed on the 'required' tag"

	err := ParseAndValidate(request, complexSpec)
	assert.DeepEqual(t, err.Error(), want)

	request.Body = ioutil.NopCloser(bytes.NewReader([]byte("{\"F\": {\"D\": true, \"E\": \"nice\", \"C\": {\"B\": 2}}}")))
	complexSpec = &complex{}

	err = ParseAndValidate(request, complexSpec)
	assert.NilError(t, err)

	request.Body = ioutil.NopCloser(bytes.NewReader([]byte("[{\"G\": \"one\"},{\"G\": \"two\"}]")))
	arrSpec := []arr{}
	want = "Key: '[0].H' Error:Field validation for 'H' failed on the 'required' tag\nKey: '[1].H' Error:Field validation for 'H' failed on the 'required' tag"

	err = ParseAndValidate(request, &arrSpec)
	assert.DeepEqual(t, err.Error(), want)
}
