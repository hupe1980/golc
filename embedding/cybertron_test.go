package embedding

import (
	"context"
	"testing"

	"github.com/nlpodyssey/cybertron/pkg/tasks/textencoding"
	"github.com/nlpodyssey/spago/mat"
	"github.com/stretchr/testify/assert"
)

// mockEncoder is a mock implementation of the textencoding.Interface interface.
type mockEncoder struct{}

func (m *mockEncoder) Encode(ctx context.Context, text string, poolingStrategy int) (textencoding.Response, error) {
	d := mat.NewDense[float32](mat.WithShape(1, 3))
	mat.SetData(d, []float32{0.1, 0.2, 0.3})

	// Mocked implementation that returns a constant embedding vector.
	return textencoding.Response{
		Vector: d,
	}, nil
}

func TestCybertron(t *testing.T) {
	t.Run("BatchEmbedText", func(t *testing.T) {
		t.Run("EmbeddingSuccess", func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			texts := []string{"text1", "text2"}
			expectedEmbeddings := [][]float32{{0.26726124, 0.5345225, 0.80178374}, {0.26726124, 0.5345225, 0.80178374}}
			encoder := &mockEncoder{}
			cybertron, err := NewCybertronFromEncoder(encoder)
			assert.NoError(t, err)

			// Act
			embeddings, err := cybertron.BatchEmbedText(ctx, texts)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expectedEmbeddings, embeddings)
		})
	})

	t.Run("EmbedText", func(t *testing.T) {
		t.Run("EmbeddingSuccess", func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			text := "text"
			expectedEmbedding := []float32{0.26726124, 0.5345225, 0.80178374}
			encoder := &mockEncoder{}
			cybertron, err := NewCybertronFromEncoder(encoder)
			assert.NoError(t, err)

			// Act
			embedding, err := cybertron.EmbedText(ctx, text)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expectedEmbedding, embedding)
		})
	})
}
