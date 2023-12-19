package documentloader

import (
	"context"
	"os"
	"strings"

	"github.com/hupe1980/golc/integration/unidoc"
	"github.com/hupe1980/golc/schema"
	"github.com/olekukonko/tablewriter"
)

// UniDocParser defines an interface for parsing documents using UniDoc.
type UniDocParser interface {
	ReadDocument(f *os.File) (unidoc.Document, error)
}

// UniDocDOCXOptions contains options for configuring the UniDocDOCX loader.
type UniDocDOCXOptions struct {
	IgnoreTables bool
}

// UniDocDOCX is a document loader for DOCX files using UniDoc.
type UniDocDOCX struct {
	parser UniDocParser
	f      *os.File
	opts   UniDocDOCXOptions
}

// NewUniDocDOCX creates a new instance of UniDocDOCX loader.
func NewUniDocDOCX(parser UniDocParser, f *os.File, optFns ...func(o *UniDocDOCXOptions)) *UniDocDOCX {
	opts := UniDocDOCXOptions{
		IgnoreTables: false,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &UniDocDOCX{
		parser: parser,
		f:      f,
		opts:   opts,
	}
}

// Load reads the document and extracts its content into schema.Document format.
func (l *UniDocDOCX) Load(ctx context.Context) ([]schema.Document, error) {
	doc, err := l.parser.ReadDocument(l.f)
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

	return []schema.Document{
		{
			PageContent: strings.Join(contents, "\n"),
			Metadata: map[string]any{
				"source": l.f.Name(),
			},
		},
	}, nil
}

// LoadAndSplit loads dOCX documents from the provided reader and splits them using the specified text splitter.
func (l *UniDocDOCX) LoadAndSplit(ctx context.Context, splitter schema.TextSplitter) ([]schema.Document, error) {
	docs, err := l.Load(ctx)
	if err != nil {
		return nil, err
	}

	return splitter.SplitDocuments(docs)
}
