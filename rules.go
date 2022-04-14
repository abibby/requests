package validate

import (
	"fmt"
)

type ValidationRule func(value any, arguments []string) error

var rules = map[string]ValidationRule{}

var initalized = false

func AddRule(key string, rule ValidationRule) {
	rules[key] = rule
}

func initRules() {
	initalized = true

	AddRule("max", func(value any, arguments []string) error {
		return fmt.Errorf("test")
	})
}

func getRule(key string) (ValidationRule, bool) {
	if !initalized {
		initRules()
	}

	r, ok := rules[key]
	return r, ok
}
