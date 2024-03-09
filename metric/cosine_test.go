package metric

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMagnitude(t *testing.T) {
	tests := []struct {
		name     string
		a        []float32
		expected float32
	}{
		{"Positive values", []float32{3, 4}, 5.0},
		{"Negative values", []float32{-3, -4}, 5.0},
		{"Mixed values", []float32{3, -4}, 5.0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Magnitude(tc.a)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCosineDistance(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []float32
		expected float32
	}{
		{"Orthogonal vectors", []float32{1, 0}, []float32{0, 1}, 1.0},
		{"Parallel vectors", []float32{1, 0}, []float32{1, 0}, 0.0},
		{"Opposite vectors", []float32{1, 0}, []float32{-1, 0}, 2.0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CosineDistance(tc.a, tc.b)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// BenchmarkCosineSimilarity benchmarks the CosineSimilarity function.
func BenchmarkCosineSimilarity(b *testing.B) {
	// Prepare random input data
	const size = 10000
	va := make([]float32, size)
	vb := make([]float32, size)

	for i := 0; i < size; i++ {
		va[i] = rand.Float32() // nolint gosec
		vb[i] = rand.Float32() // nolint gosec
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		CosineSimilarity(va, vb)
	}
}
