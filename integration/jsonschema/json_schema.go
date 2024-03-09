package jsonschema

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GenerateInline determines whether to generate inline schemas.
var GenerateInline = true

// ErrSchemaInvalid represents an error indicating an invalid schema.
var ErrSchemaInvalid = errors.New("schema is invalid")

// Mode represents the generation mode for the schema.
type Mode int

const (
	ModeAll   Mode = iota // Generate schema for all fields.
	ModeRead              // Generate schema only for fields that can be read.
	ModeWrite             // Generate schema only for fields that can be written.
)

// JSON Schema type constants
const (
	TypeBoolean = "boolean"
	TypeInteger = "integer"
	TypeNumber  = "number"
	TypeString  = "string"
	TypeArray   = "array"
	TypeObject  = "object"
)

var (
	timeType = reflect.TypeOf(time.Time{}) // timeType represents the reflection type for time.Time.
	ipType   = reflect.TypeOf(net.IP{})    // ipType represents the reflection type for net.IP.
	uriType  = reflect.TypeOf(url.URL{})   // uriType represents the reflection type for url.URL.
)

// getTagValue returns the value of a tag for the given schema and type.
func getTagValue(s *Schema, t reflect.Type, value string) (interface{}, error) {
	if s.Type == TypeString {
		return value, nil
	}

	if s.Type == TypeArray && s.Items != nil && s.Items.Type == TypeString && len(value) > 0 && value[0] != '[' {
		values := []string{}
		for _, s := range strings.Split(value, ",") {
			values = append(values, strings.TrimSpace(s))
		}

		return values, nil
	}

	var v interface{}
	if err := json.Unmarshal([]byte(value), &v); err != nil {
		return nil, err
	}

	vv := reflect.ValueOf(v)
	tv := reflect.TypeOf(v)

	if v != nil && tv != t {
		if tv.Kind() == reflect.Slice {
			tmp := reflect.MakeSlice(t, 0, vv.Len())

			for i := 0; i < vv.Len(); i++ {
				if !vv.Index(i).Elem().Type().ConvertibleTo(t.Elem()) {
					return nil, fmt.Errorf("unable to convert %v to %v: %w", vv.Index(i).Interface(), t.Elem(), ErrSchemaInvalid)
				}

				tmp = reflect.Append(tmp, vv.Index(i).Elem().Convert(t.Elem()))
			}

			v = tmp.Interface()
		} else if !tv.ConvertibleTo(t) {
			return nil, fmt.Errorf("unable to convert %v to %v: %w", tv, t, ErrSchemaInvalid)
		}

		v = reflect.ValueOf(v).Convert(t).Interface()
	}

	return v, nil
}

// Schema represents a JSON Schema which can be generated from Go structs
type Schema struct {
	Type                 string             `json:"type,omitempty"`
	Description          string             `json:"description,omitempty"`
	Items                *Schema            `json:"items,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	AdditionalProperties interface{}        `json:"additionalProperties,omitempty"`
	PatternProperties    map[string]*Schema `json:"patternProperties,omitempty"`
	Required             []string           `json:"required,omitempty"`
	Format               string             `json:"format,omitempty"`
	Enum                 []interface{}      `json:"enum,omitempty"`
	Default              interface{}        `json:"default,omitempty"`
	Example              interface{}        `json:"example,omitempty"`
	Minimum              *float64           `json:"minimum,omitempty"`
	ExclusiveMinimum     *bool              `json:"exclusiveMinimum,omitempty"`
	Maximum              *float64           `json:"maximum,omitempty"`
	ExclusiveMaximum     *bool              `json:"exclusiveMaximum,omitempty"`
	MultipleOf           float64            `json:"multipleOf,omitempty"`
	MinLength            *uint64            `json:"minLength,omitempty"`
	MaxLength            *uint64            `json:"maxLength,omitempty"`
	Pattern              string             `json:"pattern,omitempty"`
	MinItems             *uint64            `json:"minItems,omitempty"`
	MaxItems             *uint64            `json:"maxItems,omitempty"`
	UniqueItems          bool               `json:"uniqueItems,omitempty"`
	MinProperties        *uint64            `json:"minProperties,omitempty"`
	MaxProperties        *uint64            `json:"maxProperties,omitempty"`
	AllOf                []*Schema          `json:"allOf,omitempty"`
	AnyOf                []*Schema          `json:"anyOf,omitempty"`
	OneOf                []*Schema          `json:"oneOf,omitempty"`
	Not                  *Schema            `json:"not,omitempty"`
	Nullable             bool               `json:"nullable,omitempty"`
	ReadOnly             bool               `json:"readOnly,omitempty"`
	WriteOnly            bool               `json:"writeOnly,omitempty"`
	Deprecated           bool               `json:"deprecated,omitempty"`
	ContentEncoding      string             `json:"contentEncoding,omitempty"`
	Ref                  string             `json:"$ref,omitempty"`
}

// HasValidation checks if the schema has any validation rules defined.
func (s *Schema) HasValidation() bool {
	if s.Items != nil || len(s.Properties) > 0 || s.AdditionalProperties != nil || len(s.PatternProperties) > 0 || len(s.Required) > 0 || len(s.Enum) > 0 || s.Minimum != nil || s.ExclusiveMinimum != nil || s.Maximum != nil || s.ExclusiveMaximum != nil || s.MultipleOf != 0 || s.MinLength != nil || s.MaxLength != nil || s.Pattern != "" || s.MinItems != nil || s.MaxItems != nil || s.UniqueItems || s.MinProperties != nil || s.MaxProperties != nil || len(s.AllOf) > 0 || len(s.AnyOf) > 0 || len(s.OneOf) > 0 || s.Not != nil || s.Ref != "" {
		return true
	}

	return false
}

// RemoveProperty removes a property from the schema by name.
func (s *Schema) RemoveProperty(name string) {
	delete(s.Properties, name)

	for i := range s.Required {
		if s.Required[i] == name {
			s.Required[i] = s.Required[len(s.Required)-1]
			s.Required = s.Required[:len(s.Required)-1]

			break
		}
	}
}

// Generate generates a JSON Schema from the provided Go type.
func Generate(t reflect.Type) (*Schema, error) {
	return GenerateWithMode(t, ModeAll, nil, map[string]NestedSchemaReference{})
}

// GenerateFromField generates a schema from a struct field.
func GenerateFromField(f reflect.StructField, mode Mode, definedRefs map[string]NestedSchemaReference) (string, bool, *Schema, error) { // nolint gocyclo
	jsonTags := strings.Split(f.Tag.Get("json"), ",")
	name := strings.ToLower(f.Name)

	if len(jsonTags) > 0 && jsonTags[0] != "" {
		name = jsonTags[0]
	}

	if name == "-" || !f.IsExported() {
		// Skip deliberately filtered out items
		return name, false, nil, nil
	}

	s, err := GenerateWithMode(f.Type, mode, nil, definedRefs)
	if err != nil {
		return name, false, nil, err
	}

	for _, tag := range []string{"description", "doc", "format", "enum", "default", "example", "minimum", "exclusiveMinimum", "maximum", "exclusiveMaximum", "multipleOf", "minLength", "maxLength", "pattern", "minItems", "maxItems", "uniqueItems", "minProperties", "maxProperties", "nullable", "readOnly", "writeOnly", "deprecated"} {
		if tagValue, ok := f.Tag.Lookup(tag); ok {
			switch tag {
			case "description", "doc":
				s.Description = tagValue
			case "format":
				s.Format = tagValue
			case "enum":
				if err := handleEnumTag(f, s, tagValue); err != nil {
					return name, false, nil, err
				}
			case "default":
				if v, err := getTagValue(s, f.Type, tagValue); err == nil {
					s.Default = v
				} else {
					return name, false, nil, err
				}
			case "example":
				if v, err := getTagValue(s, f.Type, tagValue); err == nil {
					s.Example = v
				} else {
					return name, false, nil, err
				}
			case "minimum":
				if min, err := strconv.ParseFloat(tagValue, 64); err == nil {
					s.Minimum = &min
				} else {
					return name, false, nil, err
				}
			case "exclusiveMinimum":
				if min, err := strconv.ParseFloat(tagValue, 64); err == nil {
					s.Minimum = &min
					t := true
					s.ExclusiveMinimum = &t
				} else {
					return name, false, nil, err
				}
			case "maximum":
				if max, err := strconv.ParseFloat(tagValue, 64); err == nil {
					s.Maximum = &max
				} else {
					return name, false, nil, err
				}
			case "exclusiveMaximum":
				if max, err := strconv.ParseFloat(tagValue, 64); err == nil {
					s.Maximum = &max
					t := true
					s.ExclusiveMaximum = &t
				} else {
					return name, false, nil, err
				}
			case "multipleOf":
				if mof, err := strconv.ParseFloat(tagValue, 64); err == nil {
					s.MultipleOf = mof
				} else {
					return name, false, nil, err
				}
			case "minLength":
				if min, err := strconv.ParseUint(tagValue, 10, 64); err == nil {
					s.MinLength = &min
				} else {
					return name, false, nil, err
				}
			case "maxLength":
				if max, err := strconv.ParseUint(tagValue, 10, 64); err == nil {
					s.MaxLength = &max
				} else {
					return name, false, nil, err
				}
			case "pattern":
				s.Pattern = tagValue
				if _, err := regexp.Compile(s.Pattern); err != nil {
					return name, false, nil, err
				}
			case "minItems":
				if min, err := strconv.ParseUint(tagValue, 10, 64); err == nil {
					s.MinItems = &min
				} else {
					return name, false, nil, err
				}
			case "maxItems":
				if max, err := strconv.ParseUint(tagValue, 10, 64); err == nil {
					s.MaxItems = &max
				} else {
					return name, false, nil, err
				}
			case "uniqueItems":
				if tagValue == "true" {
					s.UniqueItems = true
				} else if tagValue == "false" {
					s.UniqueItems = false
				} else {
					return name, false, nil, fmt.Errorf("%s uniqueItems: boolean should be true or false: %w", f.Name, ErrSchemaInvalid)
				}
			case "minProperties":
				if min, err := strconv.ParseUint(tagValue, 10, 64); err == nil {
					s.MinProperties = &min
				} else {
					return name, false, nil, err
				}
			case "maxProperties":
				if max, err := strconv.ParseUint(tagValue, 10, 64); err == nil {
					s.MaxProperties = &max
				} else {
					return name, false, nil, err
				}
			case "nullable":
				if tagValue == "true" {
					s.Nullable = true
				} else if tagValue == "false" {
					s.Nullable = false
				} else {
					return name, false, nil, fmt.Errorf("%s nullable: boolean should be true or false but got %s: %w", f.Name, tagValue, ErrSchemaInvalid)
				}
			case "readOnly":
				if tagValue == "true" {
					s.ReadOnly = true
				} else if tagValue == "false" {
					s.ReadOnly = false
				} else {
					return name, false, nil, fmt.Errorf("%s readOnly: boolean should be true or false: %w", f.Name, ErrSchemaInvalid)
				}
			case "writeOnly":
				if tagValue == "true" {
					s.WriteOnly = true
				} else if tagValue == "false" {
					s.WriteOnly = false
				} else {
					return name, false, nil, fmt.Errorf("%s writeOnly: boolean should be true or false: %w", f.Name, ErrSchemaInvalid)
				}
			case "deprecated":
				if tagValue == "true" {
					s.Deprecated = true
				} else if tagValue == "false" {
					s.Deprecated = false
				} else {
					return name, false, nil, fmt.Errorf("%s deprecated: boolean should be true or false: %w", f.Name, ErrSchemaInvalid)
				}
			}
		}
	}

	optional := false

	for _, tag := range jsonTags[1:] {
		if tag == "omitempty" {
			optional = true
		}
	}

	return name, optional, s, nil
}

// handleEnumTag handles the "enum" tag for a struct field.
func handleEnumTag(f reflect.StructField, s *Schema, tagValue string) error {
	s.Enum = []interface{}{}
	enumType := f.Type
	enumSchema := s

	if s.Type == TypeArray {
		enumType = f.Type.Elem()
		enumSchema = s.Items
	}

	for _, v := range strings.Split(tagValue, ",") {
		parsed, err := getTagValue(enumSchema, enumType, v)
		if err != nil {
			return err
		}

		enumSchema.Enum = append(enumSchema.Enum, parsed)
	}

	return nil
}

type NestedSchemaReference struct {
	Name string
	Ref  string
	Type reflect.Type
}

// GenerateWithMode generates a JSON Schema with the specified mode and additional options.
func GenerateWithMode(t reflect.Type, mode Mode, schema *Schema, definedRefs map[string]NestedSchemaReference) (*Schema, error) { // nolint gocyclo
	if schema == nil {
		schema = &Schema{}
	}

	if t == ipType {
		// Special case: IP address.
		return &Schema{Type: TypeString, Format: "ipv4"}, nil
	}

	switch t.Kind() { // nolint exhaustive
	case reflect.Struct:
		// Handle special cases.
		switch t {
		case timeType:
			return &Schema{Type: TypeString, Format: "date-time"}, nil
		case uriType:
			return &Schema{Type: TypeString, Format: "uri"}, nil
		}

		if !GenerateInline {
			tname := t.Name()
			if tname != "" {
				ref, exists := definedRefs[tname]
				if exists {
					return &Schema{Ref: ref.Ref}, nil
				}

				definedRefs[tname] = NestedSchemaReference{
					Name: tname,
					Ref:  fmt.Sprintf("#/components/schemas/%s", tname),
					Type: t,
				}
			}
		}

		properties := make(map[string]*Schema)
		required := make([]string, 0)
		schema.Type = TypeObject
		schema.AdditionalProperties = false

		for _, f := range getFields(t) {
			name, optional, s, err := GenerateFromField(f, mode, definedRefs)
			if err != nil {
				return nil, err
			}

			if s == nil {
				// Skip deliberately filtered out items
				continue
			}

			if _, ok := properties[name]; ok {
				// Item already exists, ignore it since we process embedded fields
				// after top-level ones.
				continue
			}

			if s.ReadOnly && mode == ModeWrite {
				continue
			}

			if s.WriteOnly && mode == ModeRead {
				continue
			}

			properties[name] = s

			if !optional {
				required = append(required, name)
			}
		}

		if len(properties) > 0 {
			schema.Properties = properties
		}

		if len(required) > 0 {
			schema.Required = required
		}

	case reflect.Map:
		schema.Type = TypeObject

		s, err := GenerateWithMode(t.Elem(), mode, nil, definedRefs)
		if err != nil {
			return nil, err
		}

		schema.AdditionalProperties = s
	case reflect.Slice, reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			schema.Type = TypeString
		} else {
			schema.Type = TypeArray

			s, err := GenerateWithMode(t.Elem(), mode, nil, definedRefs)
			if err != nil {
				return nil, err
			}

			schema.Items = s
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		schema.Type = TypeInteger
		schema.Format = "int32"
	case reflect.Int64:
		schema.Type = TypeInteger
		schema.Format = "int64"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		schema.Type = TypeInteger
		schema.Format = "int32"
		schema.Minimum = ptr(0.0) // Unsigned integers can't be negative.
	case reflect.Uint64:
		schema.Type = TypeInteger
		schema.Format = "int64"
		schema.Minimum = ptr(0.0) // Unsigned integers can't be negative.
	case reflect.Float32:
		schema.Type = TypeNumber
		schema.Format = "float"
	case reflect.Float64:
		schema.Type = TypeNumber
		schema.Format = "double"
	case reflect.Bool:
		schema.Type = TypeBoolean
	case reflect.String:
		schema.Type = TypeString
	case reflect.Ptr:
		return GenerateWithMode(t.Elem(), mode, schema, definedRefs)
	case reflect.Interface:
		// Interfaces can be any type.
	case reflect.Uintptr, reflect.UnsafePointer, reflect.Func:
		// Ignored...
	default:
		return nil, fmt.Errorf("unsupported type %s from %s", t.Kind(), t)
	}

	return schema, nil
}

// getFields retrieves the fields of a struct type, including embedded fields recursively.
func getFields(typ reflect.Type) []reflect.StructField {
	fields := make([]reflect.StructField, 0, typ.NumField())
	embedded := []reflect.StructField{}

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Anonymous {
			embedded = append(embedded, f)
			continue
		}

		fields = append(fields, f)
	}

	for _, f := range embedded {
		newTyp := f.Type
		if newTyp.Kind() == reflect.Ptr {
			newTyp = newTyp.Elem()
		}

		if newTyp.Kind() == reflect.Struct {
			fields = append(fields, getFields(newTyp)...)
		}
	}

	return fields
}

// ptr returns a pointer to the given value.
func ptr[T any](v T) *T {
	return &v
}
