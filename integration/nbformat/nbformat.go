// Package nbformat provides utilities to read and parse Jupyter Notebook (nbformat) files.
package nbformat

import (
	"encoding/json"
	"io"
)

// Notebook represents a Jupyter Notebook containing multiple cells.
type Notebook struct {
	Cells []Cell `json:"cells"`
}

// Cell represents a single cell within a Jupyter Notebook.
type Cell struct {
	CellType string   `json:"cell_type"`
	Source   string   `json:"source"`
	Outputs  []Output `json:"outputs"`
}

// Output represents the output of a cell in a Jupyter Notebook.
type Output struct {
	ErrorName  string   `json:"ename"`
	ErrorValue string   `json:"evalue"`
	Traceback  []string `json:"traceback"`
	OutputType string   `json:"output_type"`
	Text       []string `json:"text"`
}

// ReadNBFormat reads and parses a Jupyter Notebook from the given io.Reader.
// It returns a pointer to the Notebook struct containing the parsed content.
// If there is an error during reading or parsing, it returns an error.
func ReadNBFormat(r io.Reader) (*Notebook, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var notebook Notebook
	if err = json.Unmarshal(data, &notebook); err != nil {
		return nil, err
	}

	return &notebook, nil
}
