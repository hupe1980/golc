// Package llm provides functionalities for working with Large Language Models (LLMs).
package llm

import (
	"regexp"
	"strings"
)

// EnforceStopTokens cuts off the text as soon as any stop words occur.
func EnforceStopTokens(text string, stop []string) string {
	if len(stop) == 0 {
		return text
	}

	// Create a regular expression pattern by joining stop words with "|"
	pattern := strings.Join(stop, "|")

	// Compile the regular expression pattern
	re := regexp.MustCompile(pattern)

	// Split the text using the regular expression and return the first part
	parts := re.Split(text, 2)

	return parts[0]
}
