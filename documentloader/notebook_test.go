package documentloader

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotebook(t *testing.T) {
	notebookJSON := `
	{
		"cells": [
			{
				"cell_type": "code",
				"source": "print('Hello, world!')",
				"outputs": []
			},
			{
				"cell_type": "markdown",
				"source": "# This is a Markdown cell",
				"outputs": []
			},
			{
				"cell_type": "code",
				"source": "print(10/0)",
				"outputs": [
					{
						"output_type": "error",
						"ename": "ZeroDivisionError",
						"evalue": "division by zero",
						"traceback": [
							"Traceback (most recent call last):",
							"File \"<stdin>\", line 1, in <module>",
							"ZeroDivisionError: division by zero"
						]
					}
				]
			},
			{
				"cell_type": "code",
				"source": "x = 'a' * 100",
				"outputs": [
					{
						"output_type": "stream",
						"name": "stdout",
						"text": [
							"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
						]
					}
				]
			}
		]
	}`

	reader := strings.NewReader(notebookJSON)

	notebook := &Notebook{
		r: reader,
		opts: NotebookOptions{
			IncludeOutputs:  true,
			Traceback:       true,
			MaxOutputLength: 10,
		},
	}

	documents, err := notebook.Load(context.Background())
	require.NoError(t, err, "Error loading notebook")
	require.Len(t, documents, 1, "Expected 1 document")

	expectedPageContent := "'code' cell: 'print('Hello, world!')'\n\n'markdown' cell: '# This is a Markdown cell'\n\n" +
		"'code' cell: 'print(10/0)'\n, gives error 'ZeroDivisionError', with description 'division by zero'\nand traceback '[Traceback (most recent call last): File \"<stdin>\", line 1, in <module> ZeroDivisionError: division by zero]'\n\n" +
		"'code' cell: 'x = 'a' * 100'\n with output 'aaaaaaaaaa'\n\n"
	require.Equal(t, expectedPageContent, documents[0].PageContent, "Incorrect page content")
}
