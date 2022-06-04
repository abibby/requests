package rules

import (
	"log"
	"net/http"
	"reflect"
	"strconv"
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
	Name      string
}

type ValidationRule func(options *ValidationOptions) bool

var rules = map[string]ValidationRule{}
var initalized = false

func AddRule(key string, rule ValidationRule) {
	rules[key] = rule
}

func initRules() {
	initalized = true

	initNumericRules()
	initStringRules()
	initGenericRules()
}

func GetRule(key string) (ValidationRule, bool) {
	if !initalized {
		initRules()
	}

	r, ok := rules[key]
	return r, ok
}

type TypeRuleArguments []string

func (a TypeRuleArguments) GetString(i int) string {
	if len(a) < i {
		return ""
	}
	return a[i]
}

func (a TypeRuleArguments) GetInt(i int) int64 {
	s := a.GetString(i)
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

func (a TypeRuleArguments) GetUint(i int) uint64 {
	s := a.GetString(i)
	val, _ := strconv.ParseUint(s, 10, 64)
	return val
}

func (a TypeRuleArguments) GetFloat(i int) float64 {
	s := a.GetString(i)
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

type TypeRule struct {
	ArgCount int
	Int      func(value int64, arguments TypeRuleArguments) bool
	Uint     func(value uint64, arguments TypeRuleArguments) bool
	Float    func(value float64, arguments TypeRuleArguments) bool
	String   func(value string, arguments TypeRuleArguments) bool
}

func AddTypeRule(key string, rule *TypeRule) {
	AddRule(key, func(options *ValidationOptions) bool {
		if len(options.Arguments) < rule.ArgCount {
			log.Printf("max must have %d argument(s)", rule.ArgCount)
			return true
		}
		val := reflect.ValueOf(options.Value)

		switch options.Value.(type) {
		case int, int8, int16, int32, int64:
			if rule.Int == nil {
				log.Printf("no rule for int fields")
				return true
			}
			return rule.Int(val.Int(), options.Arguments)
		case uint, uint8, uint16, uint32, uint64:
			if rule.Uint == nil {
				log.Printf("no rule for uint fields")
				return true
			}
			return rule.Uint(val.Uint(), options.Arguments)
		case float32, float64:
			if rule.Float == nil {
				log.Printf("no rule for float fields")
				return true
			}
			return rule.Float(val.Float(), options.Arguments)
		case string:
			if rule.String == nil {
				log.Printf("no rule for string fields")
				return true
			}
			return rule.String(val.String(), options.Arguments)
		default:
			log.Printf("using a numeric rule on a non numeric field")
			return true
		}
	})
}
