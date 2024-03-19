package metric

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMagnitude(t *testing.T) {
	tests := []struct {
		name     string
		vector   []float32
		expected float32
	}{
		{
			name:     "Empty Vector",
			vector:   []float32{},
			expected: 0,
		},
		{
			name:     "Single Element",
			vector:   []float32{3},
			expected: 3,
		},
		{
			name:     "Multiple Elements",
			vector:   []float32{3, 4},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Magnitude(tt.vector))
		})
	}
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name       string
		vector1    []float32
		vector2    []float32
		expected   float32
		shouldFail bool
	}{
		{
			name:       "Empty Vectors",
			vector1:    []float32{},
			vector2:    []float32{},
			expected:   0,
			shouldFail: false,
		},
		{
			name:       "Orthogonal Vectors",
			vector1:    []float32{1, 0},
			vector2:    []float32{0, 1},
			expected:   0,
			shouldFail: false,
		},
		{
			name:       "Parallel Vectors",
			vector1:    []float32{3, 4},
			vector2:    []float32{6, 8},
			expected:   1,
			shouldFail: false,
		},
		{
			name:       "Opposite Direction Vectors",
			vector1:    []float32{3, 4},
			vector2:    []float32{-3, -4},
			expected:   -1,
			shouldFail: false,
		},
		{
			name:       "Non-zero Starting Index",
			vector1:    []float32{0, 0, 1, 0},
			vector2:    []float32{0, 0, 0, 1},
			expected:   0,
			shouldFail: false,
		},
		{
			name:       "Different Length Vectors",
			vector1:    []float32{1, 2, 3},
			vector2:    []float32{4, 5, 6, 7},
			expected:   0,
			shouldFail: true,
		},
		{
			name:       "Negative Values",
			vector1:    []float32{-1, -2, -3},
			vector2:    []float32{4, 5, 6},
			expected:   -0.9746319, // Expected cosine similarity value for these vectors
			shouldFail: false,
		},
		{
			name:       "One Zero Vector",
			vector1:    []float32{0, 0, 0},
			vector2:    []float32{1, 2, 3},
			expected:   0, // One zero vector results in zero cosine similarity
			shouldFail: false,
		},
		{
			name:       "Both Zero Vectors",
			vector1:    []float32{0, 0, 0},
			vector2:    []float32{0, 0, 0},
			expected:   0, // Both zero vectors result in zero cosine similarity
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := CosineSimilarity(tt.vector1, tt.vector2)
			if tt.shouldFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

func TestSquaredL2(t *testing.T) {
	tests := []struct {
		name       string
		vector1    []float32
		vector2    []float32
		expected   float32
		shouldFail bool
	}{
		{
			name:       "Empty Vectors",
			vector1:    []float32{},
			vector2:    []float32{},
			expected:   0,
			shouldFail: false,
		},
		{
			name:       "Orthogonal Vectors",
			vector1:    []float32{1, 0},
			vector2:    []float32{0, 1},
			expected:   2,
			shouldFail: false,
		},
		{
			name:       "Identical Vectors",
			vector1:    []float32{3, 4},
			vector2:    []float32{3, 4},
			expected:   0,
			shouldFail: false,
		},
		{
			name:       "Different Length Vectors",
			vector1:    []float32{1, 2},
			vector2:    []float32{1, 2, 3},
			expected:   0,
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := SquaredL2(tt.vector1, tt.vector2)
			if tt.shouldFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			}
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
		_, _ = CosineSimilarity(va, vb)
	}
}
