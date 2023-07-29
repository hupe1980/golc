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

type HTMLOptions struct {
	TagFilter []string
}

type HTML struct {
	r    io.Reader
	opts HTMLOptions
}

func NewHTML(r io.Reader) *HTML {
	opts := HTMLOptions{
		TagFilter: []string{"script"},
	}

	return &HTML{
		r:    r,
		opts: opts,
	}
}

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
	return nil, nil
}
