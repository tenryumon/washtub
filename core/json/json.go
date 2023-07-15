package json

import (
	"reflect"
	"strings"
)

// https://stackoverflow.com/questions/40864840/how-to-get-the-json-field-names-of-a-struct-in-golang
// https://stackoverflow.com/questions/24337145/get-name-of-struct-field-using-reflection
func GetJSONField(Struct interface{}, StructField interface{}) (result string) {
	s := reflect.ValueOf(Struct).Elem()
	f := reflect.ValueOf(StructField).Elem()

	// expected: all input field have json tag on struct declaration
	for i := 0; i < s.NumField(); i++ {
		// compare the reflect value
		if s.Field(i) == f {
			targetField := s.Type().Field(i)
			jsonTag := targetField.Tag.Get("json")
			parts := strings.Split(jsonTag, ",")
			name := parts[0]
			result = name
			return
		}
	}
	return
}
