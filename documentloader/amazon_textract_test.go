package documentloader

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/hupe1980/go-textractor"
	"github.com/stretchr/testify/assert"
)

func TestAmazonTextract(t *testing.T) {
	t.Run("Load", func(t *testing.T) {
		f, err := os.Open("testdata/textract-output.json")
		assert.NoError(t, err)

		defer f.Close()

		data, err := io.ReadAll(f)
		assert.NoError(t, err)

		output := new(textractor.DocumentAPIOutput)
		err = json.Unmarshal(data, output)
		assert.NoError(t, err)

		loader := NewAmazonTextractFromOutput(output, func(o *AmazonTextractOptions) {
			o.TableLinearizationFormat = "markdown"
		})

		docs, err := loader.Load(context.Background())
		assert.NoError(t, err)

		assert.Len(t, docs, 1)
		assert.Equal(t, 1, docs[0].Metadata["page"])
		assert.Equal(t, `# New Document
## Paragraph 1
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.


| A  |  B  | C  |
|----|-----|----|
| A1 | b1  | C1 |
| A2 | B2  | C2 |
| A3 | BC3 |    |
| A4 | B4  | C4 |

`, docs[0].PageContent)
	})
}
