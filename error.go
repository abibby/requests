package validate

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gorilla/schema"
)

type ValidationError map[string][]error

var _ error = ValidationError{}

func (e ValidationError) Error() string {
	return ""
}

func (e ValidationError) HasErrors() bool {
	return len(e) > 0
}

func (e ValidationError) AddError(key string, err error) {
	if err == nil {
		return
	}
	errs := e[key]
	if errs == nil {
		errs = []error{err}
	} else {
		errs = append(errs, err)
	}
	e[key] = errs
}

func fromSchemaMultiError(err schema.MultiError) ValidationError {
	validationErr := ValidationError{}
	for key, subErr := range err {
		if err, ok := subErr.(schema.ConversionError); ok {
			validationErr[key] = []error{
				fmt.Errorf("should be of type %s", err.Type.String()),
			}
		} else {
			validationErr[key] = []error{subErr}
		}
	}
	return validationErr
}

func fromJsonUnmarshalTypeError(err *json.UnmarshalTypeError, requestStruct any) ValidationError {
	validationErr := ValidationError{}
	t := reflect.TypeOf(requestStruct)
	key := err.Field
	f, ok := t.Elem().FieldByName(err.Field)
	if ok {
		jsonKey := f.Tag.Get("json")
		if jsonKey != "" {
			key = jsonKey
		}
	}
	validationErr[key] = []error{
		fmt.Errorf("should be of type %s", err.Type.String()),
	}
	return validationErr
}
