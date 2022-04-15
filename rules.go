package validate

import (
	"fmt"
	"net/http"
)

// type Numeric interface {
// 	uint8
// 	uint16
// 	uint32
// 	uint64
// 	int8
// 	int16
// 	int32
// 	int64
// 	float32
// 	float64
// }

type ValidationOptions struct {
	Value     any
	Arguments []string
	Request   *http.Request
}

type ValidationRule func(options *ValidationOptions) error

var rules = map[string]ValidationRule{}
var initalized = false

func AddRule(key string, rule ValidationRule) {
	rules[key] = rule
}

func initRules() {
	initalized = true

	AddRule("max", func(options *ValidationOptions) error {

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
