package documentloader

import (
	"context"
	"os"
	"testing"

	"github.com/hupe1980/golc/integration/unidoc"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/textsplitter"
	"github.com/stretchr/testify/assert"
	"github.com/unidoc/unioffice/document"
)

func TestUniDocDOCX(t *testing.T) {
	t.Run("Load", func(t *testing.T) {
		tests := []struct {
			name           string
			parserMock     *uniDocParserMock
			options        UniDocDOCXOptions
			expectedResult []schema.Document
			expectedError  error
		}{
			{
				name: "Successfully load document without tables",
				parserMock: &uniDocParserMock{
					doc: getMockDocument(),
				},
				options: UniDocDOCXOptions{
					IgnoreTables: false,
				},
				expectedResult: []schema.Document{
					{
						PageContent: "Text1 from document\n\nText2 from document",
						Metadata: map[string]interface{}{
							"source": "testdata/testfile.docx",
						},
					},
				},
				expectedError: nil,
			},
			{
				name: "Successfully load document with tables",
				parserMock: &uniDocParserMock{
					doc: getMockDocumentWithTables(),
				},
				options: UniDocDOCXOptions{
					IgnoreTables: false,
				},
				expectedResult: []schema.Document{
					{
						PageContent: "Text from document\n+--------+-------+\n| Table1 | Data1 |\n+--------+-------+\n| Row1   | Value |\n+--------+-------+\n",
						Metadata: map[string]interface{}{
							"source": "testdata/testfile.docx",
						},
					},
				},
				expectedError: nil,
			},
			{
				name: "Successfully load document with ignore tables",
				parserMock: &uniDocParserMock{
					doc: getMockDocumentWithTables(),
				},
				options: UniDocDOCXOptions{
					IgnoreTables: true,
				},
				expectedResult: []schema.Document{
					{
						PageContent: "Text from document",
						Metadata: map[string]interface{}{
							"source": "testdata/testfile.docx",
						},
					},
				},
				expectedError: nil,
			},
		}

		f, err := os.Open("testdata/testfile.docx")
		assert.NoError(t, err)

		defer f.Close()

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				loader := NewUniDocDOCX(test.parserMock, f, func(o *UniDocDOCXOptions) {
					o.IgnoreTables = test.options.IgnoreTables
				})

				result, err := loader.Load(context.Background())
				assert.Equal(t, test.expectedError, err)
				assert.Equal(t, test.expectedResult, result)
			})
		}
	})

	t.Run("LoadAndSplit", func(t *testing.T) {
		tests := []struct {
			name           string
			parserMock     *uniDocParserMock
			options        UniDocDOCXOptions
			splitterMock   schema.TextSplitter
			expectedResult []schema.Document
			expectedError  error
		}{
			{
				name: "Successfully load and split document",
				parserMock: &uniDocParserMock{
					doc: getMockDocument(),
				},
				options: UniDocDOCXOptions{
					IgnoreTables: false,
				},
				splitterMock: textsplitter.NewCharacterTextSplitter(func(o *textsplitter.CharacterTextSplitterOptions) {
					o.ChunkSize = 20
				}),
				expectedResult: []schema.Document{
					{
						PageContent: "Text1 from document",
						Metadata: map[string]interface{}{
							"source": "testdata/testfile.docx",
						},
					},
					{
						PageContent: "Text2 from document",
						Metadata: map[string]interface{}{
							"source": "testdata/testfile.docx",
						},
					},
				},
				expectedError: nil,
			},
		}

		f, err := os.Open("testdata/testfile.docx")
		assert.NoError(t, err)

		defer f.Close()

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				loader := NewUniDocDOCX(test.parserMock, f, func(o *UniDocDOCXOptions) {
					o.IgnoreTables = test.options.IgnoreTables
				})

				result, err := loader.LoadAndSplit(context.Background(), test.splitterMock)
				assert.Equal(t, test.expectedError, err)
				assert.Equal(t, test.expectedResult, result)
			})
		}
	})
}

type documentMock struct {
	docText *document.DocText
}

func (m *documentMock) ExtractText() *document.DocText {
	return m.docText
}

type uniDocParserMock struct {
	doc       unidoc.Document
	parserErr error
}

func (m *uniDocParserMock) ReadDocument(f *os.File) (unidoc.Document, error) {
	if m.parserErr != nil {
		return nil, m.parserErr
	}

	return m.doc, nil
}

// Helper function to create a mock document with tables.
func getMockDocumentWithTables() unidoc.Document {
	return &documentMock{
		docText: &document.DocText{
			Items: []document.TextItem{
				{
					Text: "Text from document",
				},
				{
					TableInfo: &document.TableInfo{
						RowIndex: 0,
					},
					Text: "Table1",
				},
				{
					TableInfo: &document.TableInfo{
						RowIndex: 0,
					},
					Text: "Data1",
				},
				{
					TableInfo: &document.TableInfo{
						RowIndex: 1,
					},
					Text: "Row1",
				},
				{
					TableInfo: &document.TableInfo{
						RowIndex: 1,
					},
					Text: "Value",
				},
			},
		},
	}
}

// Helper function to create a mock document.
func getMockDocument() unidoc.Document {
	return &documentMock{
		docText: &document.DocText{
			Items: []document.TextItem{
				{
					Text: "Text1 from document\n\nText2 from document",
				},
			},
		},
	}
}
