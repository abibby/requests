package rules

import (
	"log"
	"math"
	"reflect"
	"strconv"
)

type NumberRule struct {
	ArgCount int
	Int      func(value int64, arguments []int64) bool
	Uint     func(value uint64, arguments []uint64) bool
	Float    func(value float64, arguments []float64) bool
}

func initNumericRules() {
	AddRule("max", func(options *ValidationOptions) bool {
		return AddNumericRule(options, &NumberRule{
			ArgCount: 1,
			Int: func(value int64, arguments []int64) bool {
				return value <= arguments[0]
			},
			Uint: func(value uint64, arguments []uint64) bool {
				return value <= arguments[0]
			},
			Float: func(value float64, arguments []float64) bool {
				return value <= arguments[0]
			},
		})
	})
	AddRule("min", func(options *ValidationOptions) bool {
		return AddNumericRule(options, &NumberRule{
			ArgCount: 1,
			Int: func(value int64, arguments []int64) bool {
				return value >= arguments[0]
			},
			Uint: func(value uint64, arguments []uint64) bool {
				return value >= arguments[0]
			},
			Float: func(value float64, arguments []float64) bool {
				return value >= arguments[0]
			},
		})
	})
	AddRule("multiple_of", func(options *ValidationOptions) bool {
		return AddNumericRule(options, &NumberRule{
			ArgCount: 1,
			Int: func(value int64, arguments []int64) bool {
				return value%arguments[0] == 0
			},
			Uint: func(value uint64, arguments []uint64) bool {
				return value%arguments[0] == 0
			},
			Float: func(value float64, arguments []float64) bool {
				n := value / arguments[0]
				return n == math.Floor(n)
			},
		})
	})
}

func AddNumericRule(options *ValidationOptions, rule *NumberRule) bool {
	if len(options.Arguments) < rule.ArgCount {
		log.Printf("max must have %d argument(s)", rule.ArgCount)
		return true
	}
	val := reflect.ValueOf(options.Value)

	switch options.Value.(type) {
	case int, int8, int16, int32, int64:
		intArgs := make([]int64, len(options.Arguments))
		for i, arg := range options.Arguments {
			intArg, err := strconv.ParseInt(arg, 10, 64)
			if err != nil {
				log.Printf("argument '%s' is not a valid int", options.Arguments[0])
				return true
			}
			intArgs[i] = intArg
		}
		return rule.Int(val.Int(), intArgs)
	case uint, uint8, uint16, uint32, uint64:
		uintArgs := make([]uint64, len(options.Arguments))
		for i, arg := range options.Arguments {
			uintArg, err := strconv.ParseUint(arg, 10, 64)
			if err != nil {
				log.Printf("argument '%s' is not a valid uint", options.Arguments[0])
				return true
			}
			uintArgs[i] = uintArg
		}
		return rule.Uint(val.Uint(), uintArgs)
	case float32, float64:
		floatArgs := make([]float64, len(options.Arguments))
		for i, arg := range options.Arguments {
			floatArg, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				log.Printf("argument '%s' is not a valid float", options.Arguments[0])
				return true
			}
			floatArgs[i] = floatArg
		}
		return rule.Float(val.Float(), floatArgs)
	default:
		log.Printf("using a numeric rule on a non numeric field")
		return true
	}
}
