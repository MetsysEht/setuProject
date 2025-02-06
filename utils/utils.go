package utils

import (
	"reflect"
	"strings"
)

// IsEmpty will check for given data is empty as per the go documentation
func IsEmpty(val interface{}) bool {
	if val == nil {
		return true
	}

	reflectVal := reflect.ValueOf(val)

	switch reflectVal.Kind() {
	case reflect.Int:
		return val.(int) == 0
	case reflect.Int64:
		return val.(int64) == 0
	case reflect.String:
		return strings.TrimSpace(val.(string)) == ""
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		return reflectVal.IsNil() || reflectVal.Len() == 0
	case reflect.Ptr:
		return reflectVal.IsNil()
	case reflect.Array:
		return reflectVal.IsZero()
	default:
		return false
	}
}
