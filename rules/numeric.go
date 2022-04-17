package rules

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
)

func initNumericRules() {
	AddRule("max", func(options *ValidationOptions) error {
		if len(options.Arguments) < 1 {
			return fmt.Errorf("max must have 1 argument")
		}
		val := reflect.ValueOf(options.Value)

		maxErr := fmt.Errorf("should be below max %s", options.Arguments[0])

		switch options.Value.(type) {
		case int, int8, int16, int32, int64:
			max, err := strconv.ParseInt(options.Arguments[0], 10, 64)
			if err != nil {
				log.Printf("argument '%s' is nota valid int", options.Arguments[0])
				return nil
			}
			if val.Int() > max {
				return maxErr
			}
		case uint, uint8, uint16, uint32, uint64:
			max, err := strconv.ParseUint(options.Arguments[0], 10, 64)
			if err != nil {
				log.Printf("argument '%s' is nota valid uint", options.Arguments[0])
				return nil
			}
			if val.Uint() > max {
				return maxErr
			}
		case float32, float64:
			max, err := strconv.ParseFloat(options.Arguments[0], 64)
			if err != nil {
				log.Printf("argument '%s' is nota valid float", options.Arguments[0])
				return nil
			}
			if val.Float() > max {
				return maxErr
			}
		}
		return nil
	})
}

// Digits
// Digits Between
// Max
// Min
// Multiple Of
// Between
