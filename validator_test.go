package validate

import (
	"fmt"
	"testing"

	"github.com/abibby/validate/rules"
	"github.com/stretchr/testify/assert"
)

func Test_Validate_fails_with_non_struct_arguments(t *testing.T) {
	err := Validate(nil, []string{}, 1)

	assert.Error(t, err)
}

func Test_Validate_generates_errors_on_failing_rules(t *testing.T) {
	fail := fmt.Errorf("this test always fails")
	rules.AddRule("should_fail", func(*rules.ValidationOptions) error {
		return fail
	})

	type Request struct {
		Foo int `validate:"should_fail"`
	}

	err := Validate(nil, []string{"Foo"}, &Request{})

	assert.Equal(t, ValidationError{
		"Foo": []error{fail},
	}, err)
}

func Test_Validate_ignores_failing_rules_if_no_value_is_passed(t *testing.T) {
	fail := fmt.Errorf("this test always fails")
	rules.AddRule("should_fail", func(*rules.ValidationOptions) error {
		return fail
	})

	type Request struct {
		Foo int `validate:"should_fail"`
	}

	err := Validate(nil, []string{}, &Request{})

	assert.NoError(t, err)
}

func Test_Validate_generates_no_errors_on_passing_rules(t *testing.T) {
	rules.AddRule("should_pass", func(*rules.ValidationOptions) error {
		return nil
	})

	type Request struct {
		Foo int `validate:"should_pass"`
	}

	err := Validate(nil, []string{"Foo"}, &Request{})

	assert.NoError(t, err)
}

func Test_Validate_multiple_errors(t *testing.T) {
	fail1 := fmt.Errorf("this test always fails 1")
	fail2 := fmt.Errorf("this test always fails 2")
	rules.AddRule("should_fail_1", func(*rules.ValidationOptions) error {
		return fail1
	})
	rules.AddRule("should_fail_2", func(*rules.ValidationOptions) error {
		return fail2
	})

	type Request struct {
		Foo int `validate:"should_fail_1"`
		Bar int `validate:"should_fail_1|should_fail_2"`
	}

	err := Validate(nil, []string{"Foo", "Bar"}, &Request{})

	assert.Equal(t, ValidationError{
		"Foo": []error{fail1},
		"Bar": []error{fail1, fail2},
	}, err)
}

func Test_Validate_uses_args(t *testing.T) {
	rules.AddRule("has_args", func(options *rules.ValidationOptions) error {
		return fmt.Errorf("this test uses args %s", options.Arguments[0])
	})

	type Request struct {
		Foo int `validate:"has_args:1"`
		Bar int `validate:"has_args:2"`
	}

	err := Validate(nil, []string{"Foo", "Bar"}, &Request{})

	assert.Equal(t, ValidationError{
		"Foo": []error{fmt.Errorf("this test uses args 1")},
		"Bar": []error{fmt.Errorf("this test uses args 2")},
	}, err)
}
func Test_Validate_required(t *testing.T) {
	type Request struct {
		Foo int `validate:"required"`
	}

	err := Validate(nil, []string{}, &Request{})

	assert.Equal(t, ValidationError{
		"Foo": []error{fmt.Errorf("required")},
	}, err)
}
