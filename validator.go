package validate

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func Validate(request *http.Request, keys []string, v any) error {
	s, err := getStruct(reflect.ValueOf(v))
	if err != nil {
		return errors.Wrap(err, "Validate mast take a struct or pointer to a struct")
	}
	t := s.Type()

	vErr := ValidationError{}

	for i := 0; i < s.NumField(); i++ {
		ft := t.Field(i)
		fv := s.Field(i)
		name := getName(ft)
		err := validateField(name, request, keys, ft, fv)
		if err != nil {
			vErr.Merge(err)
		}
	}

	if vErr.HasErrors() {
		return vErr
	}

	return nil
}

func validateField(name string, request *http.Request, keys []string, ft reflect.StructField, fv reflect.Value) ValidationError {
	validate, ok := ft.Tag.Lookup("validate")
	if !ok {
		return nil
	}

	vErr := ValidationError{}

	rulesStr := strings.Split(validate, "|")
	for _, ruleStr := range rulesStr {
		ruleName, argsStr := split(ruleStr, ":")
		args := strings.Split(argsStr, ",")

		hasKey := includes(keys, name)
		if !hasKey {
			if ruleName == "required" {
				vErr.AddError(name, fmt.Errorf("required"))
			} else {
				return nil
			}
		} else {
			rule, ok := getRule(ruleName)
			if !ok {
				continue
			}

			err := rule(&ValidationOptions{
				Value:     fv.Interface(),
				Arguments: args,
				Request:   request,
			})
			if err != nil {
				vErr.AddError(name, err)
			}

		}
	}
	return vErr
}

func getName(f reflect.StructField) string {
	name, ok := f.Tag.Lookup("json")
	if ok {
		return name
	}

	name, ok = f.Tag.Lookup("query")
	if ok {
		return name
	}

	return f.Name
}

func getStruct(v reflect.Value) (reflect.Value, error) {
	switch v.Kind() {
	case reflect.Struct:
		return v, nil
	case reflect.Interface, reflect.Pointer:
		return getStruct(v.Elem())
	default:
		return reflect.Value{}, fmt.Errorf("expected struct found %s", v.Kind())
	}
}

func split(s, sep string) (string, string) {
	parts := strings.SplitN(s, sep, 2)
	switch len(parts) {
	case 0:
		return "", ""
	case 1:
		return parts[0], ""
	default:
		return parts[0], parts[1]
	}
}

func includes[T comparable](haystack []T, needle T) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
