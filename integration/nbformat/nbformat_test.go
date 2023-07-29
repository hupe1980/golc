package nbformat

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadNBFormat(t *testing.T) {
	ipynbContent := `
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
			}
		]
	}`

	reader := strings.NewReader(ipynbContent)

	notebook, err := ReadNBFormat(reader)
	require.NoError(t, err, "Error reading nbformat data")

	expectedNumCells := 2
	require.Len(t, notebook.Cells, expectedNumCells)

	require.Equal(t, "print('Hello, world!')", notebook.Cells[0].Source)
	require.Equal(t, "code", notebook.Cells[0].CellType)
	require.Equal(t, "# This is a Markdown cell", notebook.Cells[1].Source)
	require.Equal(t, "markdown", notebook.Cells[1].CellType)
}
