package documentloader

import (
	"context"
	"strings"
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestCSV(t *testing.T) {
	t.Run("TestLoad", func(t *testing.T) {
		csvData := `id,name,age,gender
1,John,30,Male
2,Alice,25,Female
3,Bob,35,Male`

		r := strings.NewReader(csvData)
		loader := NewCSV(r)

		expectedLoad := []schema.Document{
			{
				PageContent: "id: 1\nname: John\nage: 30\ngender: Male",
				Metadata:    map[string]interface{}{"row": uint(1)},
			},
			{
				PageContent: "id: 2\nname: Alice\nage: 25\ngender: Female",
				Metadata:    map[string]interface{}{"row": uint(2)},
			},
			{
				PageContent: "id: 3\nname: Bob\nage: 35\ngender: Male",
				Metadata:    map[string]interface{}{"row": uint(3)},
			},
		}

		docsLoad, err := loader.Load(context.Background())
		assert.NoError(t, err)
		assert.ElementsMatch(t, expectedLoad, docsLoad)
	})

	t.Run("TestLoadWithFilter", func(t *testing.T) {
		csvData := `id,name,age,gender
1,John,30,Male
2,Alice,25,Female
3,Bob,35,Male`

		r := strings.NewReader(csvData)
		loader := NewCSV(r, "name", "age")

		expectedLoadWithFilter := []schema.Document{
			{
				PageContent: "name: John\nage: 30",
				Metadata:    map[string]interface{}{"row": uint(1)},
			},
			{
				PageContent: "name: Alice\nage: 25",
				Metadata:    map[string]interface{}{"row": uint(2)},
			},
			{
				PageContent: "name: Bob\nage: 35",
				Metadata:    map[string]interface{}{"row": uint(3)},
			},
		}

		docsLoadWithFilter, err := loader.Load(context.Background())
		assert.NoError(t, err)
		assert.ElementsMatch(t, expectedLoadWithFilter, docsLoadWithFilter)
	})
}
