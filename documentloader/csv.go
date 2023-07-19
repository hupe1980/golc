package documentloader

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
)

// Compile time check to ensure CSV satisfies the DocumentLoader interface.
var _ schema.DocumentLoader = (*CSV)(nil)

// CSV represents a CSV document loader.
type CSV struct {
	r       io.Reader
	columns []string
}

// NewCSV creates a new CSV loader with an io.Reader and optional column names for filtering.
func NewCSV(r io.Reader, columns ...string) *CSV {
	return &CSV{
		r:       r,
		columns: columns,
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
			if len(l.columns) > 0 && !util.Contains(l.columns, header[i]) {
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
