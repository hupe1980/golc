package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHTMLAndGetStrippedStrings(t *testing.T) {
	tests := []struct {
		name        string
		htmlContent string
		expected    string
		expectedErr error
	}{
		{
			name:        "Valid HTML",
			htmlContent: `<html><body><p>Hello, <b>world</b>!</p></body></html>`,
			expected:    "Hello, world!",
			expectedErr: nil,
		},
		{
			name:        "Empty HTML",
			htmlContent: "",
			expected:    "",
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := ParseHTMLAndGetStrippedStrings(test.htmlContent)
			assert.Equal(t, test.expected, result)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
