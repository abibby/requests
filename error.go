package validate

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gorilla/schema"
)

type ValidationError map[string][]string

var _ error = ValidationError{}

func (e ValidationError) Error() string {
	return ""
}

func fromSchemaMultiError(err schema.MultiError) ValidationError {
	validationErr := ValidationError{}
	for key, subErr := range err {
		if err, ok := subErr.(schema.ConversionError); ok {
			validationErr[key] = []string{
				fmt.Sprintf("should be of type %s", err.Type.String()),
			}
		} else {
			validationErr[key] = []string{subErr.Error()}
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
	validationErr[key] = []string{
		fmt.Sprintf("should be of type %s", err.Type.String()),
	}
	return validationErr
}
