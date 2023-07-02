package util

import (
	"reflect"
	"strings"
)

func StructToMap(obj interface{}) map[string]interface{} {
	objValue := reflect.ValueOf(obj)
	objType := objValue.Type()

	result := make(map[string]interface{})

	for i := 0; i < objValue.NumField(); i++ {
		field := objValue.Field(i)
		fieldType := objType.Field(i)

		tag := fieldType.Tag.Get("map")

		// Skip fields with "-" or "omitempty" tag
		if tag == "-" || (strings.Contains(tag, "omitempty") && fieldIsEmpty(field)) {
			continue
		}

		var value any

		if fieldType.Type.Kind() == reflect.Struct {
			fieldMap := StructToMap(field.Interface())
			value = fieldMap
		} else {
			value = field.Interface()
		}

		key := fieldType.Name

		if tag != "" {
			tagParts := strings.Split(tag, ",")
			if len(tagParts) > 0 {
				key = tagParts[0] // Use the custom tag as the key
			}
		}

		result[key] = value
	}

	return result
}

// fieldIsEmpty checks if a field is empty or its value is the zero value for the field's type.
func fieldIsEmpty(field reflect.Value) bool {
	switch field.Kind() { // nolint exhaustive
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return field.Len() == 0
	case reflect.Bool:
		return !field.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return field.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return field.IsNil()
	default:
		zeroValue := reflect.Zero(field.Type()).Interface()
		currentValue := field.Interface()

		return reflect.DeepEqual(currentValue, zeroValue)
	}
}
