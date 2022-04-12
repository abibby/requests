package validate

import (
	"fmt"
	"os"
	"reflect"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

func Validate(v any) error {
	s, err := getStruct(reflect.ValueOf(v))
	if err != nil {
		return errors.Wrap(err, "Validate mast take a struct or pointer to a struct")
	}
	t := s.Type()

	for i := 0; i < s.NumField(); i++ {
		ft := t.Field(i)
		fv := s.Field(i)
		validate, ok := ft.Tag.Lookup("validate")
		spew.Dump(ok, validate)
		spew.Dump(fv.Interface())
	}

	spew.Dump(s)
	os.Exit(1)
	return nil
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
