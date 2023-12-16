// Package deepcopy provides a function for creating deep copies of values.
package deepcopy

import (
	"reflect"
	"time"
)

// Copy creates a deep copy of the provided source and returns it as an interface{}.
// The returned value may need to be asserted to the correct type.
func Copy(src interface{}) interface{} {
	if src == nil {
		return nil
	}

	// Make the interface a reflect.Value.
	v := reflect.ValueOf(src)

	// Make a copy of the same type as the original.
	c := reflect.New(v.Type()).Elem()

	// Recursively copy the original.
	copyRecursive(v, c)

	// Return the copy as an interface.
	return c.Interface()
}

// copyRecursive performs the actual copying of the interface.
func copyRecursive(o, c reflect.Value) {
	switch o.Kind() {
	case reflect.Ptr:
		// Get the actual value being pointed to.
		ov := o.Elem()

		// If it isn't valid, return.
		if !ov.IsValid() {
			return
		}

		// Set the copied value to a new instance of the same type.
		c.Set(reflect.New(ov.Type()))

		// Recursively copy the pointed-to value.
		copyRecursive(ov, c.Elem())

	case reflect.Interface:
		// If this is a nil interface, don't do anything.
		if o.IsNil() {
			return
		}

		// Get the value for the interface, not the pointer.
		ov := o.Elem()

		// Create a new value by calling Elem().
		copyValue := reflect.New(ov.Type()).Elem()

		// Recursively copy the interface value.
		copyRecursive(ov, copyValue)
		c.Set(copyValue)

	case reflect.Struct:
		// Check if the struct field is of type time.Time and handle it accordingly.
		t, ok := o.Interface().(time.Time)
		if ok {
			c.Set(reflect.ValueOf(t))
			return
		}

		// Go through each field of the struct and copy it.
		for i := 0; i < o.NumField(); i++ {
			// The Type's StructField for a given field is checked to see if StructField.PkgPath
			// is set to determine if the field is exported or not because CanSet() returns false
			// for settable fields.
			if o.Type().Field(i).PkgPath != "" {
				continue
			}

			// Recursively copy the struct field.
			copyRecursive(o.Field(i), c.Field(i))
		}

	case reflect.Slice:
		// If the slice is nil, don't do anything.
		if o.IsNil() {
			return
		}

		// Make a new slice and copy each element.
		c.Set(reflect.MakeSlice(o.Type(), o.Len(), o.Cap()))

		for i := 0; i < o.Len(); i++ {
			copyRecursive(o.Index(i), c.Index(i))
		}

	case reflect.Map:
		// If the map is nil, don't do anything.
		if o.IsNil() {
			return
		}

		// Make a new map and copy each key-value pair.
		c.Set(reflect.MakeMap(o.Type()))

		for _, key := range o.MapKeys() {
			ov := o.MapIndex(key)
			copyValue := reflect.New(ov.Type()).Elem()

			// Recursively copy the map value.
			copyRecursive(ov, copyValue)

			// Copy the key and set the map entry.
			copyKey := Copy(key.Interface())
			c.SetMapIndex(reflect.ValueOf(copyKey), copyValue)
		}

	default:
		// For all other types, simply set the value.
		c.Set(o)
	}
}
