package validate

import (
	"fmt"
)

const (
	// Error Message for Integer Validator
	MessageNonZero   = "%s tidak boleh 0."
	MessageMinNumber = "%s harus lebih besar daripada %d."
	MessageMaxNumber = "%s harus lebih kecil daripada %d."
)

type IntegerOption func(varname string, value int64) (bool, string)

var listCommonIntegerValidation = map[string][]IntegerOption{}

func Integer(varname string, value int64, options ...IntegerOption) []string {
	newValue := int64(value)
	errors := []string{}
	for _, function := range options {
		if valid, message := function(varname, newValue); !valid {
			errors = append(errors, message)
		}

	}
	return errors
}

func AddIntegerCommon(key string, options ...IntegerOption) {
	listCommonIntegerValidation[key] = options
}

func GetIntegerCommon(key string) []IntegerOption {
	return listCommonIntegerValidation[key]
}

func WithNonZero() IntegerOption {
	return func(varname string, value int64) (bool, string) {
		if value == 0 {
			return false, fmt.Sprintf(MessageNonZero, varname)
		}
		return true, ""
	}
}

func WithPositive() IntegerOption {
	return func(varname string, value int64) (bool, string) {
		if value <= 0 {
			return false, fmt.Sprintf(MessageMinNumber, varname, 0)
		}
		return true, ""
	}
}

func WithMinimum(number int64) IntegerOption {
	return func(varname string, value int64) (bool, string) {
		if value < number {
			return false, fmt.Sprintf(MessageMinNumber, varname, number)
		}
		return true, ""
	}
}

func WithMaximum(number int64) IntegerOption {
	return func(varname string, value int64) (bool, string) {
		if value > number {
			return false, fmt.Sprintf(MessageMaxNumber, varname, number)
		}
		return true, ""
	}
}
