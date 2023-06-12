package textsplitter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitText(t *testing.T) {
	t.Run("with overlap", func(t *testing.T) {
		text := "foo bar baz 123"
		splitter := NewCharacterTextSplitter(func(o *CharacterTextSplitterOptions) {
			o.Separator = " "
			o.ChunkSize = 7
			o.ChunkOverlap = 3
		})

		chunks := splitter.splitText(text)

		assert.ElementsMatch(t, chunks, []string{"foo bar", "bar baz", "baz 123"})
	})

	t.Run("ignores empty docs", func(t *testing.T) {
		text := "foo  bar"
		splitter := NewCharacterTextSplitter(func(o *CharacterTextSplitterOptions) {
			o.Separator = " "
			o.ChunkSize = 2
			o.ChunkOverlap = 0
		})

		chunks := splitter.splitText(text)

		assert.ElementsMatch(t, chunks, []string{"foo", "bar"})
	})
}
