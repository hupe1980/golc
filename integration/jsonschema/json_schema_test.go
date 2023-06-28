package jsonschema

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	t.Run("Generate schema for struct with required and optional fields", func(t *testing.T) {
		type MyStruct struct {
			RequiredField int    `json:"required_field"`
			OptionalField string `json:"optional_field,omitempty"`
		}

		schema, err := Generate(reflect.TypeOf(MyStruct{}))
		assert.NoError(t, err)

		// Assert that the generated schema has the expected properties
		assert.Equal(t, 2, len(schema.Properties))
		assert.Contains(t, schema.Properties, "required_field")
		assert.Contains(t, schema.Properties, "optional_field")
	})

	t.Run("Generate schema for struct with enum field", func(t *testing.T) {
		type MyEnum string

		const (
			EnumValue1 MyEnum = "value1"
			EnumValue2 MyEnum = "value2"
		)

		type MyStruct struct {
			EnumField MyEnum `json:"enum_field" enum:"value1,value2"`
		}

		schema, err := Generate(reflect.TypeOf(MyStruct{}))
		assert.NoError(t, err)

		// Assert that the generated schema has the expected enum values
		assert.Equal(t, 2, len(schema.Properties["enum_field"].Enum))
		assert.Contains(t, schema.Properties["enum_field"].Enum, string(EnumValue1))
		assert.Contains(t, schema.Properties["enum_field"].Enum, string(EnumValue2))
	})

	t.Run("Generate schema for struct with nested structs", func(t *testing.T) {
		type NestedStruct struct {
			NestedField int `json:"nested_field"`
		}

		type MyStruct struct {
			Nested NestedStruct `json:"nested"`
		}

		schema, err := Generate(reflect.TypeOf(MyStruct{}))
		assert.NoError(t, err)

		// Assert that the generated schema has the expected properties
		assert.Equal(t, 1, len(schema.Properties))
		assert.Contains(t, schema.Properties, "nested")
		assert.Equal(t, 1, len(schema.Properties["nested"].Properties))
		assert.Contains(t, schema.Properties["nested"].Properties, "nested_field")
	})

	t.Run("Generate schema for struct with array field", func(t *testing.T) {
		type MyStruct struct {
			ArrayField []string `json:"array_field"`
		}

		schema, err := Generate(reflect.TypeOf(MyStruct{}))
		assert.NoError(t, err)

		// Assert that the generated schema has the expected properties
		assert.Equal(t, 1, len(schema.Properties))
		assert.Contains(t, schema.Properties, "array_field")
		assert.Equal(t, "array", schema.Properties["array_field"].Type)
	})

	t.Run("Generate schema for struct with ignored fields", func(t *testing.T) {
		type MyStruct struct {
			VisibleField int    `json:"visible_field"`
			IgnoredField string `json:"-"`
		}

		schema, err := Generate(reflect.TypeOf(MyStruct{}))
		assert.NoError(t, err)

		// Assert that the generated schema has the expected properties
		assert.Equal(t, 1, len(schema.Properties))
		assert.Contains(t, schema.Properties, "visible_field")
		assert.NotContains(t, schema.Properties, "ignored_field")
	})

	t.Run("Generate schema for struct with pointer field", func(t *testing.T) {
		type MyStruct struct {
			PointerField *string `json:"pointer_field"`
		}

		schema, err := Generate(reflect.TypeOf(MyStruct{}))
		assert.NoError(t, err)

		// Assert that the generated schema has the expected properties
		assert.Equal(t, 1, len(schema.Properties))
		assert.Contains(t, schema.Properties, "pointer_field")
	})

	t.Run("Generate schema for struct with time.Time field", func(t *testing.T) {
		type MyStruct struct {
			TimeField time.Time `json:"time_field"`
		}

		schema, err := Generate(reflect.TypeOf(MyStruct{}))
		assert.NoError(t, err)

		// Assert that the generated schema has the expected properties
		assert.Equal(t, 1, len(schema.Properties))
		assert.Contains(t, schema.Properties, "time_field")
		assert.Equal(t, "string", schema.Properties["time_field"].Type)
		assert.Equal(t, "date-time", schema.Properties["time_field"].Format)
	})

	t.Run("Generate schema for struct with map field", func(t *testing.T) {
		type MyStruct struct {
			MapField map[string]int `json:"map_field"`
		}

		schema, err := Generate(reflect.TypeOf(MyStruct{}))
		assert.NoError(t, err)

		// Assert that the generated schema has the expected properties
		assert.Equal(t, 1, len(schema.Properties))
		assert.Contains(t, schema.Properties, "map_field")
		assert.Equal(t, "object", schema.Properties["map_field"].Type)
	})

	t.Run("Generate schema for struct with anonymous field", func(t *testing.T) {
		type AnonymousStruct struct {
			AnonymousField int `json:"anonymous_field"`
		}

		type MyStruct struct {
			AnonymousStruct
			VisibleField string `json:"visible_field"`
		}

		schema, err := Generate(reflect.TypeOf(MyStruct{}))
		assert.NoError(t, err)

		// Assert that the generated schema has the expected properties
		assert.Equal(t, 2, len(schema.Properties))
		assert.Contains(t, schema.Properties, "anonymous_field")
		assert.Contains(t, schema.Properties, "visible_field")
	})

	t.Run("Generate schema for struct with custom JSON tags", func(t *testing.T) {
		type MyStruct struct {
			FieldOne   int    `json:"field_one,omitempty"`
			FieldTwo   string `json:"field_two,omitempty"`
			FieldThree bool   `json:"field_three,omitempty"`
		}

		schema, err := Generate(reflect.TypeOf(MyStruct{}))
		assert.NoError(t, err)

		// Assert that the generated schema has the expected properties
		assert.Equal(t, 3, len(schema.Properties))
		assert.Contains(t, schema.Properties, "field_one")
		assert.Contains(t, schema.Properties, "field_two")
		assert.Contains(t, schema.Properties, "field_three")
	})

	t.Run("Generate schema for struct with unsupported field types", func(t *testing.T) {
		type MyStruct struct {
			Channel chan int `json:"channel"`
			Func    func()   `json:"func"`
		}

		_, err := Generate(reflect.TypeOf(MyStruct{}))
		assert.Error(t, err)
	})

	t.Run("Generate schema for struct with private fields", func(t *testing.T) {
		type MyStruct struct {
			PublicField  int `json:"public_field"`
			privateField int //nolint unused
		}

		schema, err := Generate(reflect.TypeOf(MyStruct{}))
		assert.NoError(t, err)

		// Assert that the generated schema has the expected properties
		assert.Equal(t, 1, len(schema.Properties))
		assert.Contains(t, schema.Properties, "public_field")
	})
}
