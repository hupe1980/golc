package sqldb

import (
	"regexp"
	"strings"
)

// CleanQuery cleans sql query from double white space, comments and leading/trailing spaces.
func CleanQuery(query string) string {
	// remove comments
	query = regexp.MustCompile(`/\*(.*)\*/|\-\-(.*)`).ReplaceAllString(query, "")

	// remove double white space
	query = regexp.MustCompile(`\s+`).ReplaceAllString(query, " ")

	// remove leading/trailing space
	query = strings.TrimSpace(query)

	return query
}
