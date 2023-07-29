package documentloader

import (
	"context"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure HTML satisfies the DocumentLoader interface.
var _ schema.DocumentLoader = (*HTML)(nil)

// HTMLOptions contains options for the HTML document loader.
type HTMLOptions struct {
	// TagFilter is a list of HTML tags to be filtered from the document content.
	TagFilter []string
}

// HTML implements the DocumentLoader interface for HTML documents.
type HTML struct {
	r    io.Reader
	opts HTMLOptions
}

// NewHTML creates a new HTML document loader with an io.Reader and optional configuration options.
// It returns a pointer to the created HTML loader.
func NewHTML(r io.Reader, optFns ...func(o *HTMLOptions)) *HTML {
	opts := HTMLOptions{
		TagFilter: []string{"script"},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &HTML{
		r:    r,
		opts: opts,
	}
}

// Load loads the HTML document from the reader and extracts the text content.
// It returns a list of schema.Document containing the extracted content and the title as metadata.
func (l *HTML) Load(ctx context.Context) ([]schema.Document, error) {
	doc, err := goquery.NewDocumentFromReader(l.r)
	if err != nil {
		return nil, err
	}

	title := doc.Find("title").Text()

	var textSlice []string

	var sel *goquery.Selection
	if doc.Has("body") != nil {
		sel = doc.Find("body").Contents()
	} else {
		sel = doc.Contents()
	}

	sel.Each(func(_ int, s *goquery.Selection) {
		for _, ft := range l.opts.TagFilter {
			if s.Is(ft) {
				return
			}
		}

		text := strings.TrimSpace(s.Text())
		if text != "" {
			textSlice = append(textSlice, text)
		}
	})

	return []schema.Document{
		{
			PageContent: strings.Join(textSlice, "\n"),
			Metadata: map[string]any{
				"title": title,
			},
		},
	}, nil
}

// LoadAndSplit loads HTML documents from the provided reader and splits them using the specified text splitter.
func (l *HTML) LoadAndSplit(ctx context.Context, splitter schema.TextSplitter) ([]schema.Document, error) {
	docs, err := l.Load(ctx)
	if err != nil {
		return nil, err
	}

	return splitter.SplitDocuments(docs)
}
