package rules

import (
	"testing"
)

func TestString(t *testing.T) {
	data := map[string]TestCase{
		"alpha-pass":      {"alpha", &ValidationOptions{Value: "a"}, true},
		"alpha-fail":      {"alpha", &ValidationOptions{Value: "1a"}, false},
		"alpha_dash-pass": {"alpha_dash", &ValidationOptions{Value: "a-"}, true},
		"alpha_dash-fail": {"alpha_dash", &ValidationOptions{Value: "1a-"}, false},
		"alpha_num-pass":  {"alpha_num", &ValidationOptions{Value: "a1"}, true},
		"alpha_num-fail":  {"alpha_num", &ValidationOptions{Value: "1a!"}, false},
		"numeric-pass":    {"numeric", &ValidationOptions{Value: "123"}, true},
		"numeric-fail":    {"numeric", &ValidationOptions{Value: "123a"}, false},
		"email-pass":      {"email", &ValidationOptions{Value: "user@example.com"}, true},
		"email-fail":      {"email", &ValidationOptions{Value: "not an email"}, false},
		"ends_with-pass":  {"ends_with", &ValidationOptions{Value: "stringend", Arguments: []string{"end"}}, true},
		"ends_with-fail":  {"ends_with", &ValidationOptions{Value: "stringendnot", Arguments: []string{"end"}}, false},
	}

	runTests(t, data)
}
