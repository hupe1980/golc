package util

import "strings"

// Capitalize capitalizes the first letter of the given string.
// If the input string is empty, it returns the input string itself.
func Capitalize(s string) string {
	if s == "" {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}
