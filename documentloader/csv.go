package documentloader

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure CSV satisfies the DocumentLoader interface.
var _ schema.DocumentLoader = (*CSV)(nil)

// CSVOptions contains options for configuring the CSV loader.
type CSVOptions struct {
	// Separator is the rune used to separate fields in the CSV file.
	Separator rune

	// LazyQuotes controls whether the CSV reader should use lazy quotes mode.
	LazyQuotes bool

	// Columns is a list of column names to filter and include in the loaded documents.
	Columns []string
}

// CSV represents a CSV document loader.
type CSV struct {
	r    io.Reader
	opts CSVOptions
}

// NewCSV creates a new CSV loader with an io.Reader and optional configuration options.
// It returns a pointer to the created CSV loader.
func NewCSV(r io.Reader, optFns ...func(o *CSVOptions)) *CSV {
	opts := CSVOptions{
		Separator:  ',',
		LazyQuotes: false,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &CSV{
		r:    r,
		opts: opts,
	}
}

// Load loads CSV documents from the provided reader.
func (l *CSV) Load(ctx context.Context) ([]schema.Document, error) {
	var (
		header []string
		docs   []schema.Document
		rown   uint
	)

	reader := csv.NewReader(l.r)
	reader.Comma = l.opts.Separator
	reader.LazyQuotes = l.opts.LazyQuotes

	isHeader := true

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if isHeader {
			header = row
			isHeader = false

			continue
		}

		var content []string

		for i, value := range row {
			if len(l.opts.Columns) > 0 && !util.Contains(l.opts.Columns, header[i]) {
				continue
			}

			line := fmt.Sprintf("%s: %s", header[i], value)
			content = append(content, line)
		}

		rown++
		docs = append(docs, schema.Document{
			PageContent: strings.Join(content, "\n"),
			Metadata:    map[string]any{"row": rown},
		})
	}

	return docs, nil
}

// LoadAndSplit loads CSV documents from the provided reader and splits them using the specified text splitter.
func (l *CSV) LoadAndSplit(ctx context.Context, splitter schema.TextSplitter) ([]schema.Document, error) {
	docs, err := l.Load(ctx)
	if err != nil {
		return nil, err
	}

	return splitter.SplitDocuments(docs)
}
