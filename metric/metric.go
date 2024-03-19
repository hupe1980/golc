package metric

import (
	"errors"

	"github.com/hupe1980/golc/internal/math32"
)

// Magnitude calculates the magnitude (length) of a float32 slice.
func Magnitude(v []float32) float32 {
	return math32.Sqrt(math32.Dot(v, v))
}

// CosineSimilarity calculates the cosine similarity between two float32 slices.
func CosineSimilarity(v1, v2 []float32) (float32, error) {
	// Check if the vector sizes match
	if len(v1) != len(v2) {
		return 0, errors.New("vector sizes do not match")
	}

	dotProduct := math32.Dot(v1, v2)
	sumA := math32.Dot(v1, v1)
	sumB := math32.Dot(v2, v2)

	// Avoid division by zero
	if sumA == 0 || sumB == 0 {
		return 0, nil
	}

	return dotProduct / math32.Sqrt(sumA*sumB), nil
}

// CosineDistance calculates the cosine distance between two float32 slices.
func CosineDistance(a, b []float32) (float32, error) {
	s, err := CosineSimilarity(a, b)
	if err != nil {
		return 0, err
	}

	return 1 - s, nil
}

// SquaredL2 calculates the squared L2 distance between two float32 slices.
func SquaredL2(v1, v2 []float32) (float32, error) {
	// Check if the vector sizes match
	if len(v1) != len(v2) {
		return 0, errors.New("vector sizes do not match")
	}

	return math32.SquaredL2(v1, v2), nil
}
