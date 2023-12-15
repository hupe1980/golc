package util

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ParseHTMLAndGetStrippedStrings parses the HTML content and returns the stripped strings.
// It uses the goquery package to extract the text from HTML elements.
func ParseHTMLAndGetStrippedStrings(htmlContent string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var strippedStrings []string

	doc.Find("body > *").Each(func(_ int, s *goquery.Selection) {
		strippedString := strings.TrimSpace(s.Text())
		if strippedString != "" {
			strippedStrings = append(strippedStrings, strippedString)
		}
	})

	return strings.Join(strippedStrings, " "), nil
}
