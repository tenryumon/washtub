package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// ref: https://articles.wesionary.team/reflections-tutorial-query-string-to-struct-parser-in-go-b2f858f99ea1
// QueryParser maps get request into struct
func QueryParser(req *http.Request, d interface{}) error {
	dType := reflect.TypeOf(d)

	if err := shouldBeStruct(dType); err != nil {
		return err
	}

	// Data Holder Value
	dhVal := reflect.ValueOf(d)

	// Loop over all the fields present in struct (Title, Body, JSON)
	for i := 0; i < dType.Elem().NumField(); i++ {

		// Give me ith field. Elem() is used to dereference the pointer
		field := dType.Elem().Field(i)

		// Get the value from field tag i.e in case of Title it is "title"
		key := field.Tag.Get("json")
		if key == "" {
			continue
		}

		// Get the type of field
		kind := field.Type.Kind()

		// Get the value from query params with given key
		val := req.URL.Query().Get(key)

		//  Get reference of field value provided to input `d`
		result := dhVal.Elem().Field(i)

		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if result.CanSet() && val != "" {
				intVal, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return fmt.Errorf("failed to parsed %s, err %v", key, err)
				}
				if result.CanSet() {
					result.SetInt(intVal)
				}
			}
			break
		case reflect.String:
			if result.CanSet() {
				result.SetString(val)
			}
			break
		case reflect.Bool:
			if result.CanSet() {
				result.SetBool(strings.ToLower(val) == "true")
			}
			break
		case reflect.Slice, reflect.Struct:
			if result.CanSet() && val != "" {
				valType := reflect.New(field.Type)
				err := json.Unmarshal([]byte(val), valType.Interface())
				if err != nil {
					return fmt.Errorf("failed to parsed %s, err %v", key, err)
				}
				result.Set(valType.Elem())
			}
			break
		}
	}
	return nil
}

func shouldBeStruct(d reflect.Type) error {
	td := d.Elem()
	if td.Kind() != reflect.Struct {
		errStr := fmt.Sprintf("Input should be %v, found %v", reflect.Struct, td.Kind())
		return errors.New(errStr)
	}
	return nil
}
