// Package nbformat provides utilities to read and parse Jupyter Notebook (nbformat) files.
package nbformat

import (
	"encoding/json"
	"io"
	"strings"
)

// Notebook represents a Jupyter Notebook containing multiple cells.
type Notebook struct {
	Metadata      Metadata `json:"metadata"`
	Nbformat      int      `json:"nbformat"`
	NbformatMinor int      `json:"nbformat_minor"`
	Cells         []Cell   `json:"cells"`
}

// Metadata represents the metadata of a Jupyter Notebook.
type Metadata struct {
	KernelSpec   KernelSpec   `json:"kernelspec"`
	LanguageInfo LanguageInfo `json:"language_info"`
}

// KernelSpec represents the kernel specification in the metadata.
type KernelSpec struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

// LanguageInfo represents the language information in the metadata.
type LanguageInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Cell represents a single cell within a Jupyter Notebook.
type Cell struct {
	CellType string                 `json:"cell_type"`
	Source   string                 `json:"source"` // Could be []string, but we always convert it to a single string
	Metadata map[string]interface{} `json:"metadata"`
	Outputs  []Output               `json:"outputs,omitempty"`
}

// UnmarshalJSON custom unmarshals a Cell to ensure Source is always a single string.
func (c *Cell) UnmarshalJSON(data []byte) error {
	type Alias Cell

	aux := &struct {
		Source interface{} `json:"source"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch v := aux.Source.(type) {
	case string:
		c.Source = v
	case []interface{}:
		var lines []string
		for _, line := range v {
			lines = append(lines, line.(string))
		}

		c.Source = strings.Join(lines, "\n")
	}

	return nil
}

// Output represents the output of a cell in a Jupyter Notebook.
type Output struct {
	OutputType string                 `json:"output_type"`
	Text       string                 `json:"text,omitempty"` // Could be []string, but we always convert it to a single string
	Data       map[string]interface{} `json:"data,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	ErrorName  string                 `json:"ename,omitempty"`
	ErrorValue string                 `json:"evalue,omitempty"`
	Traceback  []string               `json:"traceback,omitempty"`
}

// UnmarshalJSON custom unmarshals an Output to ensure Text is always a single string.
func (o *Output) UnmarshalJSON(data []byte) error {
	type Alias Output

	aux := &struct {
		Text interface{} `json:"text,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch v := aux.Text.(type) {
	case string:
		o.Text = v
	case []interface{}:
		var lines []string
		for _, line := range v {
			lines = append(lines, line.(string))
		}

		o.Text = strings.Join(lines, "\n")
	}

	return nil
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
