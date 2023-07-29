package documentloader

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/hupe1980/golc/integration/nbformat"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Notebook satisfies the DocumentLoader interface.
var _ schema.DocumentLoader = (*Notebook)(nil)

// NotebookOptions represents the options for loading a Jupyter Notebook.
type NotebookOptions struct {
	// Include outputs (cell execution results) in the document content.
	IncludeOutputs bool

	// Include traceback information for cells with errors.
	Traceback bool

	// Maximum length of output text to include in the document.
	MaxOutputLength uint
}

// Notebook represents a Jupyter Notebook document loader.
type Notebook struct {
	// Reader to read the notebook content.
	r io.Reader

	// Options for loading the notebook.
	opts NotebookOptions
}

// NewNotebook creates a new instance of Notebook with the given reader and optional functions to set options.
func NewNotebook(r io.Reader, optFns ...func(o *NotebookOptions)) *Notebook {
	opts := NotebookOptions{
		IncludeOutputs:  false,
		Traceback:       false,
		MaxOutputLength: 10,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Notebook{
		r:    r,
		opts: opts,
	}
}

// Load reads and parses the Jupyter Notebook from the provided reader.
// It returns a slice of schema.Document representing the notebook's content and any error encountered.
func (l *Notebook) Load(ctx context.Context) ([]schema.Document, error) {
	notebook, err := nbformat.ReadNBFormat(l.r)
	if err != nil {
		return nil, err
	}

	pageContent := ""

	for _, c := range notebook.Cells {
		if l.opts.IncludeOutputs && len(c.Outputs) > 0 {
			if c.Outputs[0].ErrorName != "" {
				eName := c.Outputs[0].ErrorName
				eValue := c.Outputs[0].ErrorValue

				if l.opts.Traceback {
					traceback := c.Outputs[0].Traceback
					pageContent += fmt.Sprintf("'%s' cell: '%s'\n, gives error '%s', with description '%s'\nand traceback '%s'\n\n", c.CellType, c.Source, eName, eValue, traceback)
				} else {
					pageContent += fmt.Sprintf("'%s' cell: '%s'\n, gives error '%s', with description '%s'\n\n", c.CellType, c.Source, eName, eValue)
				}
			} else if c.Outputs[0].OutputType == "stream" {
				output := strings.Join(c.Outputs[0].Text, "")
				minOutput := len(output)

				if minOutput > int(l.opts.MaxOutputLength) {
					minOutput = int(l.opts.MaxOutputLength)
				}

				pageContent += fmt.Sprintf("'%s' cell: '%s'\n with output '%s'\n\n", c.CellType, c.Source, output[:minOutput])
			}
		} else {
			pageContent += fmt.Sprintf("'%s' cell: '%s'\n\n", c.CellType, c.Source)
		}
	}

	return []schema.Document{
		{
			PageContent: pageContent,
			Metadata:    map[string]any{},
		},
	}, nil
}

// LoadAndSplit loads Notebook documents from the provided reader and splits them using the specified text splitter.
// It returns a slice of schema.Document representing the notebook's content and any error encountered.
func (l *Notebook) LoadAndSplit(ctx context.Context, splitter schema.TextSplitter) ([]schema.Document, error) {
	docs, err := l.Load(ctx)
	if err != nil {
		return nil, err
	}

	return splitter.SplitDocuments(docs)
}
