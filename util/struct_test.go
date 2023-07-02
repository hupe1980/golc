package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructToMap(t *testing.T) {
	// Test case 1: Empty struct
	t.Run("EmptyStruct", func(t *testing.T) {
		obj := struct{}{}
		expected := map[string]interface{}{}
		result := StructToMap(obj)
		assert.Equal(t, expected, result, "Unexpected map conversion result")
	})

	// Test case 2: Struct with non-empty fields
	t.Run("NonEmptyStruct", func(t *testing.T) {
		obj := struct {
			Name string `map:"name"`
			Age  int    `map:"age"`
		}{
			Name: "John",
			Age:  30,
		}
		expected := map[string]interface{}{
			"name": "John",
			"age":  30,
		}
		result := StructToMap(obj)
		assert.Equal(t, expected, result, "Unexpected map conversion result")
	})

	// Test case 3: Struct with empty field and omitempty tag
	t.Run("StructWithEmptyFieldAndOmitEmptyTag", func(t *testing.T) {
		obj := struct {
			Name string `map:"name,omitempty"`
			Age  int    `map:"age,omitempty"`
		}{
			Name: "",
			Age:  0,
		}
		expected := map[string]interface{}{}
		result := StructToMap(obj)
		assert.Equal(t, expected, result, "Unexpected map conversion result")
	})

	// Test case 4: Struct with nested struct
	t.Run("StructWithNestedStruct", func(t *testing.T) {
		obj := struct {
			Name    string `map:"name"`
			Address struct {
				Street string `map:"street"`
				City   string `map:"city"`
			} `map:"address"`
		}{
			Name: "John",
			Address: struct {
				Street string `map:"street"`
				City   string `map:"city"`
			}{
				Street: "123 Main St",
				City:   "New York",
			},
		}
		expected := map[string]interface{}{
			"name": "John",
			"address": map[string]interface{}{
				"street": "123 Main St",
				"city":   "New York",
			},
		}
		result := StructToMap(obj)
		assert.Equal(t, expected, result, "Unexpected map conversion result")
	})

	// Test case 5: Struct with ignored field
	t.Run("StructWithIgnoredField", func(t *testing.T) {
		obj := struct {
			Name  string `map:"name"`
			Email string `map:"-"`
			Age   int    `map:"age"`
		}{
			Name:  "John",
			Email: "john@example.com",
			Age:   30,
		}
		expected := map[string]interface{}{
			"name": "John",
			"age":  30,
		}
		result := StructToMap(obj)
		assert.Equal(t, expected, result, "Unexpected map conversion result")
	})
}
