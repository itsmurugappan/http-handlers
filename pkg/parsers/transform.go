package parsers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var ErrInvalidInput = errors.New("input must be a struct pointer")
var validate = validator.New()

func JsonToSpec(data []byte, spec interface{}) error {
	s := reflect.ValueOf(spec)
	if s.Kind() != reflect.Ptr {
		return ErrInvalidInput
	}

	t := reflect.TypeOf(spec)
	new := reflect.New(t.Elem())
	i := new.Interface()
	json.Unmarshal(data, i)

	s.Elem().Set((reflect.ValueOf(i)).Elem())

	return nil
}

func ParseReqBodyToSpec(r *http.Request, spec interface{}) error {
	body, _ := ioutil.ReadAll(r.Body)
	return JsonToSpec(body, spec)
}

func ParseAndValidate(r *http.Request, spec interface{}) error {
	if err := ParseReqBodyToSpec(r, spec); err != nil {
		return err
	}
	switch (reflect.ValueOf(spec)).Elem().Kind() {
	case reflect.Struct:
		return validate.Struct(spec)
	case reflect.Slice:
		return validate.Var(spec, "dive")
	default:
		return validate.Struct(spec)
	}
}
