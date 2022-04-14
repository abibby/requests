package validate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Validate_fails_with_non_struct_arguments(t *testing.T) {
	err := Validate(1)

	assert.Error(t, err)
}

func Test_Validate_generates_errors_on_failing_rules(t *testing.T) {
	fail := fmt.Errorf("this test always fails")
	AddRule("should_fail", func(value any, arguments []string) error {
		return fail
	})

	type Request struct {
		Foo int `validate:"should_fail:1"`
	}

	err := Validate(&Request{})

	assert.Equal(t, err, ValidationError{
		"Foo": []error{fail},
	})
}

func Test_Validate_generates_no_errors_on_passing_rules(t *testing.T) {
	AddRule("should_pass", func(value any, arguments []string) error {
		return nil
	})

	type Request struct {
		Foo int `validate:"should_pass:1"`
	}

	err := Validate(&Request{})

	assert.NoError(t, err)
}

func Test_Validate_multiple_errors(t *testing.T) {
	fail1 := fmt.Errorf("this test always fails 1")
	fail2 := fmt.Errorf("this test always fails 2")
	AddRule("should_fail_1", func(value any, arguments []string) error {
		return fail1
	})
	AddRule("should_fail_2", func(value any, arguments []string) error {
		return fail2
	})

	type Request struct {
		Foo int `validate:"should_fail_1"`
		Bar int `validate:"should_fail_1|should_fail_2"`
	}

	err := Validate(&Request{})

	assert.Equal(t, err, ValidationError{
		"Foo": []error{fail1},
		"Bar": []error{fail1, fail2},
	})
}
