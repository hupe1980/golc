package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCosineSimilarity(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		matrix1 := [][]float64{{1.0, 2.0}, {3.0, 4.0}}
		matrix2 := [][]float64{{1.0, 0.0}, {0.0, 1.0}}

		expected := 0.6454972243679028 // Expected cosine similarity value

		result := CosineSimilarity(matrix1, matrix2)
		assert.InDelta(t, expected, result, 1e-9)
	})

	t.Run("Zero matrices", func(t *testing.T) {
		matrix1 := [][]float64{}
		matrix2 := [][]float64{}

		expected := 0.0 // Both matrices are empty, so cosine similarity should be 0

		result := CosineSimilarity(matrix1, matrix2)
		assert.InDelta(t, expected, result, 1e-9)
	})

	t.Run("Matrices with all zeros", func(t *testing.T) {
		matrix1 := [][]float64{{0.0, 0.0}, {0.0, 0.0}}
		matrix2 := [][]float64{{0.0, 0.0}, {0.0, 0.0}}

		expected := 0.0 // Both matrices have all zeros, so cosine similarity should be 0

		result := CosineSimilarity(matrix1, matrix2)
		assert.InDelta(t, expected, result, 1e-9)
	})

	t.Run("Matrices with orthogonal vectors", func(t *testing.T) {
		matrix1 := [][]float64{{1.0, 0.0}, {0.0, 1.0}}
		matrix2 := [][]float64{{0.0, 1.0}, {1.0, 0.0}}

		expected := 0.0 // Matrices have orthogonal vectors, so cosine similarity should be 0

		result := CosineSimilarity(matrix1, matrix2)
		assert.InDelta(t, expected, result, 1e-9)
	})

	t.Run("Matrices with identical vectors", func(t *testing.T) {
		matrix1 := [][]float64{{1.0, 2.0}, {3.0, 4.0}}
		matrix2 := [][]float64{{1.0, 2.0}, {3.0, 4.0}}

		expected := 1.0 // Matrices have identical vectors, so cosine similarity should be 1

		result := CosineSimilarity(matrix1, matrix2)
		assert.InDelta(t, expected, result, 1e-9)
	})

}
