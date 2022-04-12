package validate

type ValidationRule func(value string, arguments []string) error

var rules map[string]ValidationRule

func AddRule(string, ValidationRule) {

}
