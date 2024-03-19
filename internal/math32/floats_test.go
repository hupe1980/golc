package math32

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDot(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []float32
		expected float32
	}{
		{"Positive values (size 3)", []float32{1, 2, 3}, []float32{4, 5, 6}, 32.0},
		{"Negative values (size 3)", []float32{-1, -2, -3}, []float32{-4, -5, -6}, 32.0},
		{"More than 4 (size 6)", []float32{1, 2, 3, 1, 2, 3}, []float32{4, 5, 6, 4, 5, 6}, 64.0},
		{"Mixed values (size 3)", []float32{1, -2, 3}, []float32{-4, 5, -6}, -32.0},
		{"Zero values (size 3)", []float32{0, 0, 0}, []float32{0, 0, 0}, 0.0},
		{"Positive values (size 9)", []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}, []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}, 285.0},
		{"Positive values (size 10)", []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 385.0},
		{"Positive values (size 15)", []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, 1240.0},
		{"Positive values (size 16)", []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, 1496.0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Dot(tc.a, tc.b)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// BenchmarkDot-10    	    7623	    157954 ns/op	       0 B/op	       0 allocs/op
func BenchmarkDot(b *testing.B) {
	// Generate random float32 slices for benchmarking.
	const size = 1000000 // Size of slices
	va := make([]float32, size)
	vb := make([]float32, size)

	for i := range va {
		va[i] = rand.Float32() // nolint gosec
		vb[i] = rand.Float32() // nolint gosec
	}

	// Run the Dot function b.N times and measure the time taken.
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Dot(va, vb)
	}
}

func TestSquaredL2(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []float32
		expected float32
	}{
		{"Positive values", []float32{1, 2, 3}, []float32{4, 5, 6}, 27.0},
		{"Negative values", []float32{-1, -2, -3}, []float32{-4, -5, -6}, 27.0},
		{"1 Remainder", []float32{1, 2, 3, 1, 2, 3}, []float32{4, 5, 6, 4, 5, 6}, 54.0},
		{"Mixed values", []float32{1, -2, 3}, []float32{-4, 5, -6}, 155.0},
		{"Zero values", []float32{0, 0, 0}, []float32{0, 0, 0}, 0.0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := SquaredL2(tc.a, tc.b)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// BenchmarkSquaredL2-10    	    5128	    235120 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSquaredL2(b *testing.B) {
	// Generate random float32 slices for benchmarking.
	const size = 1000000 // Size of slices
	va := make([]float32, size)
	vb := make([]float32, size)

	for i := range va {
		va[i] = rand.Float32() // nolint gosec
		vb[i] = rand.Float32() // nolint gosec
	}

	// Run the Dot function b.N times and measure the time taken.
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = SquaredL2(va, vb)
	}
}
