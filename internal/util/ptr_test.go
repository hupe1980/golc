package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddrOrNil(t *testing.T) {
	t.Run("ZeroValue", func(t *testing.T) {
		var zeroInt int
		result := AddrOrNil(zeroInt)
		assert.Nil(t, result, "Expected nil for zero value")
	})

	t.Run("NonZeroValue", func(t *testing.T) {
		nonZeroInt := 42
		result := AddrOrNil(nonZeroInt)
		assert.NotNil(t, result, "Expected non-nil for non-zero value")
		assert.Equal(t, nonZeroInt, *result, "Unexpected value for non-zero value")
	})

	t.Run("ZeroString", func(t *testing.T) {
		var zeroString string
		result := AddrOrNil(zeroString)
		assert.Nil(t, result, "Expected nil for zero value")
	})

	t.Run("NonZeroString", func(t *testing.T) {
		nonZeroString := "test"
		result := AddrOrNil(nonZeroString)
		assert.NotNil(t, result, "Expected non-nil for non-zero value")
		assert.Equal(t, nonZeroString, *result, "Unexpected value for non-zero value")
	})
}

func TestPTR(t *testing.T) {
	tests := []struct {
		name string
		in   interface{}
	}{
		{
			name: "Integer",
			in:   42,
		},
		{
			name: "String",
			in:   "test",
		},
		{
			name: "Struct",
			in: struct {
				Name  string
				Value int
			}{
				Name:  "example",
				Value: 123,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PTR(tt.in)

			// Check if the pointer is not nil
			assert.NotNil(t, result)

			// Check if the value pointed to is equal to the original value
			assert.Equal(t, tt.in, *result)
		})
	}
}
