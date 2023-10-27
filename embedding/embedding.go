// Package embedding contains the implementation to create vector embeddings
// from text using different APIs
package embedding

import "strings"

func removeNewLines(text string) string {
	return strings.ReplaceAll(text, "\n", " ")
}
