package documentloader

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/textsplitter"
	"github.com/stretchr/testify/assert"
)

func TestText_Load(t *testing.T) {
	t.Run("Load", func(t *testing.T) {
		tests := []struct {
			name     string
			input    io.Reader
			expected []schema.Document
			err      error
		}{
			{
				name:  "Empty Reader",
				input: bytes.NewReader([]byte("")),
				expected: []schema.Document{
					{
						PageContent: "",
						Metadata:    map[string]interface{}{},
					},
				},
				err: nil,
			},
			{
				name:  "Single Document",
				input: bytes.NewReader([]byte("This is a test document.")),
				expected: []schema.Document{
					{
						PageContent: "This is a test document.",
						Metadata:    map[string]interface{}{},
					},
				},
				err: nil,
			},
			{
				name:     "Error Reading",
				input:    failingReader{},
				expected: nil,
				err:      io.ErrUnexpectedEOF,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				loader := NewText(test.input)
				docs, err := loader.Load(context.Background())

				assert.Equal(t, test.err, err)
				assert.Equal(t, test.expected, docs)
			})
		}
	})

	t.Run("LoadAndSplit", func(t *testing.T) {
		tests := []struct {
			name     string
			input    io.Reader
			expected []schema.Document
			err      error
		}{
			{
				name:     "Empty Reader",
				input:    bytes.NewReader([]byte("")),
				expected: []schema.Document{},
				err:      nil,
			},
			{
				name:  "Single Document",
				input: bytes.NewReader([]byte("This is a test document.")),
				expected: []schema.Document{
					{
						PageContent: "This is a test document.",
						Metadata:    map[string]interface{}{},
					},
				},
				err: nil,
			},
			{
				name:     "Error Reading",
				input:    failingReader{},
				expected: nil,
				err:      io.ErrUnexpectedEOF,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				loader := NewText(test.input)
				docs, err := loader.LoadAndSplit(context.Background(), textsplitter.NewRecusiveCharacterTextSplitter())

				assert.Equal(t, test.err, err)
				assert.Equal(t, test.expected, docs)
			})
		}
	})
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}
