package validate

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
)

const (
	// Error Message for String Validator
	MessageRequired    = "%s harus diisi."
	MessageMinLength   = "%s minimal harus terdiri %d karakter."
	MessageMaxLength   = "%s maksimal harus terdiri %d karakter."
	MessageOnlyAlpha   = "%s tidak boleh mengandung angka."
	MessageOnlyNumeric = "%s harus berupa angka."
	MessageOnlyAscii   = "%s tidak boleh mengandung karakter aneh."
	MessagePhoneFormat = "%s tidak valid."
	MessageEmailFormat = "%s tidak valid."
	MessageTimeFormat  = "%s tidak valid."
)

type StringOption func(varname string, value string) (bool, string)

var listCommonStringValidation = map[string][]StringOption{}

func String(varname string, value string, options ...StringOption) []string {
	errors := []string{}
	for _, function := range options {
		if valid, message := function(varname, value); !valid {
			errors = append(errors, message)
		}

	}
	return errors
}

func AddStringCommon(key string, options ...StringOption) {
	listCommonStringValidation[key] = options
}

func GetStringCommon(key string) []StringOption {
	return listCommonStringValidation[key]
}

func WithRequired() StringOption {
	return func(varname string, value string) (bool, string) {
		if len(value) == 0 {
			return false, fmt.Sprintf(MessageRequired, varname)
		}
		return true, ""
	}
}

func WithMinLength(min int) StringOption {
	return func(varname string, value string) (bool, string) {
		if value == "" {
			// If empty, just return because maybe is optional
			return true, ""
		}

		if len(value) < min {
			return false, fmt.Sprintf(MessageMinLength, varname, min)
		}
		return true, ""
	}
}

func WithMaxLength(max int) StringOption {
	return func(varname string, value string) (bool, string) {
		if value == "" {
			// If empty, just return because maybe is optional
			return true, ""
		}

		if len(value) > max {
			return false, fmt.Sprintf(MessageMaxLength, varname, max)
		}
		return true, ""
	}
}

func WithOnlyAlpha() StringOption {
	return func(varname string, value string) (bool, string) {
		if value == "" {
			// If empty, just return because maybe is optional
			return true, ""
		}

		for _, char := range value {
			if !isCharAscii(char) {
				return false, fmt.Sprintf(MessageOnlyAlpha, varname)
			}
			if isCharNumber(char) {
				return false, fmt.Sprintf(MessageOnlyAlpha, varname)
			}
		}
		return true, ""
	}
}

func WithOnlyNumeric() StringOption {
	return func(varname string, value string) (bool, string) {
		if value == "" {
			// If empty, just return because maybe is optional
			return true, ""
		}

		for _, char := range value {
			if !isCharNumber(char) {
				return false, fmt.Sprintf(MessageOnlyNumeric, varname)
			}
		}
		return true, ""
	}
}

func WithAlphaNumeric() StringOption {
	return func(varname string, value string) (bool, string) {
		if value == "" {
			// If empty, just return because maybe is optional
			return true, ""
		}

		for _, char := range value {
			if !isCharAscii(char) {
				return false, fmt.Sprintf(MessageOnlyAscii, varname)
			}
		}
		return true, ""
	}
}

func WithPhoneFormat() StringOption {
	return func(varname string, value string) (bool, string) {
		if value == "" {
			// If empty, just return because maybe is optional
			return true, ""
		}

		if !isPhonePrefix(value) {
			return false, fmt.Sprintf(MessagePhoneFormat, varname)
		}
		if strings.HasPrefix(value, "+62") {
			value = strings.Replace(value, "+", "", 1)
		}

		for _, char := range value {
			if !isCharNumber(char) {
				return false, fmt.Sprintf(MessagePhoneFormat, varname)
			}
		}

		return true, ""
	}
}

const emailPattern = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"

var emailRegex = regexp.MustCompile(emailPattern)

func WithEmailFormat() StringOption {
	return func(varname string, value string) (bool, string) {
		if value == "" {
			// If empty, just return because maybe is optional
			return true, ""
		}

		if !emailRegex.Match([]byte(value)) {
			return false, fmt.Sprintf(MessageEmailFormat, varname)
		}

		arr := strings.Split(value, "@")
		if strings.Contains(arr[0], "+") {
			return false, fmt.Sprintf(MessageEmailFormat, varname)
		}

		return true, ""
	}
}

func WithTimeFormat(format string) StringOption {
	format = strings.ReplaceAll(format, "YYYY", "2006")
	format = strings.ReplaceAll(format, "MM", "01")
	format = strings.ReplaceAll(format, "DD", "02")
	format = strings.ReplaceAll(format, "HH", "15")
	format = strings.ReplaceAll(format, "MI", "04")
	format = strings.ReplaceAll(format, "SS", "05")

	return func(varname string, value string) (bool, string) {
		if value == "" {
			// If empty, just return because maybe is optional
			return true, ""
		}

		_, err := time.Parse(format, value)
		if err != nil {
			return false, fmt.Sprintf(MessageTimeFormat, varname)
		}
		return true, ""
	}
}

func isCharNumber(char rune) bool {
	return char >= '0' && char <= '9'
}

func isCharAscii(char rune) bool {
	return char <= unicode.MaxASCII
}

func isPhonePrefix(phone string) bool {
	return strings.HasPrefix(phone, "0") || strings.HasPrefix(phone, "62") || strings.HasPrefix(phone, "+62")
}
