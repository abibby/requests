package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Validate_fails_with_non_struct_arguments(t *testing.T) {
	err := Validate(1)

	assert.Error(t, err)
}

func Test_Validate_(t *testing.T) {
	type Request struct {
		Foo string `validate:"max:1"`
	}

	err := Validate(&Request{
		Foo: "a",
	})

	assert.Error(t, err)
}
