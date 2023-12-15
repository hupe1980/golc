package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapitalize(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Non-empty string",
			input:    "hello",
			expected: "Hello",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Single character",
			input:    "a",
			expected: "A",
		},
		{
			name:     "Already capitalized",
			input:    "World",
			expected: "World",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the Capitalize function
			result := Capitalize(tt.input)

			// Check the result against the expected value
			assert.Equal(t, tt.expected, result)
		})
	}
}
