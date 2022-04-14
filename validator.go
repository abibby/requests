package validate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func Validate(v any) error {
	s, err := getStruct(reflect.ValueOf(v))
	if err != nil {
		return errors.Wrap(err, "Validate mast take a struct or pointer to a struct")
	}
	t := s.Type()

	vErr := ValidationError{}

	for i := 0; i < s.NumField(); i++ {
		ft := t.Field(i)
		fv := s.Field(i)
		validate, ok := ft.Tag.Lookup("validate")
		if !ok {
			continue
		}

		rulesStr := strings.Split(validate, "|")
		for _, ruleStr := range rulesStr {
			ruleName, argsStr := split(ruleStr, ":")
			args := strings.Split(argsStr, ",")

			rule, ok := getRule(ruleName)
			if !ok {
				continue
			}

			err = rule(fv.Interface(), args)
			if err != nil {
				vErr.AddError(getName(ft), err)
			}
		}

		// name := getName(ft)

		fmt.Printf("%#v\n", vErr)
	}

	if vErr.HasErrors() {
		return vErr
	}

	return nil
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
