package documentloader

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTML(t *testing.T) {
	tests := []struct {
		name          string
		htmlContent   string
		expectedText  string
		expectedTitle string
	}{
		{
			name: "HTML document with body element and title",
			htmlContent: `<!DOCTYPE html>
<html>
<head>
	<title>Test Document</title>
</head>
<body>
	<h1>Hello World!</h1>
	<p>This is a test document.</p>
</body>
</html>`,
			expectedText:  "Hello World!\nThis is a test document.",
			expectedTitle: "Test Document",
		},
		{
			name: "HTML document without body element",
			htmlContent: `<h1>Hello World!</h1>
<p>This is a test document.</p>`,
			expectedText:  "Hello World!\nThis is a test document.",
			expectedTitle: "",
		},
		{
			name: "HTML document without filtered text",
			htmlContent: `<h1>Hello World!</h1>
<script>console.log("This should be filtered")</script>
<p>This is a test document.</p>`,
			expectedText:  "Hello World!\nThis is a test document.",
			expectedTitle: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.htmlContent)
			htmlLoader := NewHTML(r)
			docs, err := htmlLoader.Load(context.Background())
			require.NoError(t, err, "Unexpected error")
			require.Len(t, docs, 1, "Expected 1 document, but got %d", len(docs))
			require.Equal(t, tt.expectedText, docs[0].PageContent, "Incorrect page content")
			require.Equal(t, tt.expectedTitle, docs[0].Metadata["title"])
		})
	}
}
