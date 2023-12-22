package documentloader

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/hupe1980/golc/integration/unidoc"
	"github.com/hupe1980/golc/schema"
	"github.com/olekukonko/tablewriter"
)

// UniDocParser defines an interface for parsing documents using UniDoc.
type UniDocParser interface {
	ReadDocument(r io.ReaderAt, size int64) (unidoc.Document, error)
}

// UniDocDOCXOptions contains options for configuring the UniDocDOCX loader.
type UniDocDOCXOptions struct {
	IgnoreTables bool

	// Source is the name of the pdf document
	Source string
}

// UniDocDOCX is a document loader for DOCX files using UniDoc.
type UniDocDOCX struct {
	parser UniDocParser
	r      io.ReaderAt
	size   int64
	opts   UniDocDOCXOptions
}

func NewUniDocDOCX(parser UniDocParser, r io.ReaderAt, size int64, optFns ...func(o *UniDocDOCXOptions)) *UniDocDOCX {
	opts := UniDocDOCXOptions{
		IgnoreTables: false,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &UniDocDOCX{
		parser: parser,
		r:      r,
		size:   size,
		opts:   opts,
	}
}

// NewUniDocDOCX creates a new instance of UniDocDOCX loader.
func NewUniDocDOCXFromFile(parser UniDocParser, f *os.File, optFns ...func(o *UniDocDOCXOptions)) *UniDocDOCX {
	opts := UniDocDOCXOptions{
		IgnoreTables: false,
		Source:       f.Name(),
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return NewUniDocDOCX(parser, f, 0, func(o *UniDocDOCXOptions) {
		o.IgnoreTables = opts.IgnoreTables
		o.Source = opts.Source
	})
}

// Load reads the document and extracts its content into schema.Document format.
func (l *UniDocDOCX) Load(ctx context.Context) ([]schema.Document, error) {
	doc, err := l.parser.ReadDocument(l.r, l.size)
	if err != nil {
		return nil, err
	}

	contents := []string{}
	rowIndex := -1
	tableData := [][]string{}

	extracted := doc.ExtractText()
	for i, e := range extracted.Items {
		text := ""

		if tblInfo := e.TableInfo; tblInfo != nil {
			if l.opts.IgnoreTables {
				continue
			}

			if rowIndex != tblInfo.RowIndex { // new row
				rowIndex = tblInfo.RowIndex

				tableData = append(tableData, []string{})
			}

			tableData[tblInfo.RowIndex] = append(tableData[tblInfo.RowIndex], e.Text)
		} else {
			text = e.Text
		}

		if e.TableInfo == nil || i == len(extracted.Items)-1 {
			if len(tableData) > 0 {
				b := new(strings.Builder)
				table := tablewriter.NewWriter(b)

				table.SetRowLine(true)
				table.AppendBulk(tableData)
				table.Render()

				contents = append(contents, b.String())

				rowIndex = -1
				tableData = [][]string{}
			}
		}

		if text != "" {
			contents = append(contents, text)
		}
	}

	textDoc := schema.Document{
		PageContent: strings.Join(contents, "\n"),
	}

	if l.opts.Source != "" {
		textDoc.Metadata = map[string]any{
			"source": l.opts.Source,
		}
	}

	return []schema.Document{textDoc}, nil
}

// LoadAndSplit loads dOCX documents from the provided reader and splits them using the specified text splitter.
func (l *UniDocDOCX) LoadAndSplit(ctx context.Context, splitter schema.TextSplitter) ([]schema.Document, error) {
	docs, err := l.Load(ctx)
	if err != nil {
		return nil, err
	}

	return splitter.SplitDocuments(docs)
}
