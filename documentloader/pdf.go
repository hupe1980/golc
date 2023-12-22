package documentloader

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hupe1980/golc/schema"
	"github.com/ledongthuc/pdf"
)

// Compile time check to ensure PDF satisfies the DocumentLoader interface.
var _ schema.DocumentLoader = (*PDF)(nil)

type PDFOptions struct {
	// Password for encrypted PDF files.
	Password string

	// Page number to start loading from (default is 1).
	StartPage uint

	// Maximum number of pages to load (0 for all pages).
	MaxPages uint

	// Source is the name of the pdf document
	Source string
}

// PDF represents a PDF document loader that implements the DocumentLoader interface.
type PDF struct {
	f    io.ReaderAt
	size int64
	opts PDFOptions
}

// NewPDFFromFile creates a new PDF loader with the given options.
func NewPDF(f io.ReaderAt, size int64, optFns ...func(o *PDFOptions)) (*PDF, error) {
	opts := PDFOptions{
		StartPage: 1,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.StartPage == 0 {
		opts.StartPage = 1
	}

	return &PDF{
		f:    f,
		size: size,
		opts: opts,
	}, nil
}

// NewPDFFromFile creates a new PDF loader with the given options.
func NewPDFFromFile(f *os.File, optFns ...func(o *PDFOptions)) (*PDF, error) {
	opts := PDFOptions{
		StartPage: 1,
		Source:    f.Name(),
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.StartPage == 0 {
		opts.StartPage = 1
	}

	finfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return NewPDF(f, finfo.Size(), func(o *PDFOptions) {
		*o = opts
	})
}

// Load loads the PDF document and returns a slice of schema.Document containing the page contents and metadata.
func (l *PDF) Load(ctx context.Context) ([]schema.Document, error) {
	var (
		reader *pdf.Reader
		err    error
	)

	if l.opts.Password != "" {
		reader, err = pdf.NewReaderEncrypted(l.f, l.size, func() string {
			return l.opts.Password
		})
		if err != nil {
			return nil, err
		}
	} else {
		reader, err = pdf.NewReader(l.f, l.size)
		if err != nil {
			return nil, err
		}
	}

	numPages := reader.NumPage()
	if l.opts.StartPage > uint(numPages) {
		return nil, fmt.Errorf("startpage out of page range: 1-%d", numPages)
	}

	maxPages := numPages - int(l.opts.StartPage) + 1
	if l.opts.MaxPages > 0 && numPages > int(l.opts.MaxPages) {
		maxPages = int(l.opts.MaxPages)
	}

	docs := make([]schema.Document, 0, numPages)

	fonts := make(map[string]*pdf.Font)

	page := 1

	for i := int(l.opts.StartPage); i < maxPages+int(l.opts.StartPage); i++ {
		p := reader.Page(i)

		for _, name := range p.Fonts() {
			if _, ok := fonts[name]; !ok {
				f := p.Font(name)
				fonts[name] = &f
			}
		}

		text, err := p.GetPlainText(fonts)
		if err != nil {
			return nil, err
		}

		// add the document to the doc list
		doc := schema.Document{
			PageContent: strings.TrimSpace(text),
			Metadata: map[string]any{
				"page":       page,
				"totalPages": maxPages,
			},
		}

		if l.opts.Source != "" {
			doc.Metadata["source"] = l.opts.Source
		}

		docs = append(docs, doc)

		page++
	}

	return docs, nil
}

// LoadAndSplit loads PDF documents from the provided reader and splits them using the specified text splitter.
func (l *PDF) LoadAndSplit(ctx context.Context, splitter schema.TextSplitter) ([]schema.Document, error) {
	docs, err := l.Load(ctx)
	if err != nil {
		return nil, err
	}

	return splitter.SplitDocuments(docs)
}
