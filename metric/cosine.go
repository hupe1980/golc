package metric

import (
	"github.com/hupe1980/golc/internal/math32"
)

// Magnitude calculates the magnitude (length) of a float32 slice.
func Magnitude(a []float32) float32 {
	return math32.Sqrt(math32.Dot(a, a))
}

// CosineSimilarity calculates the cosine similarity between two float32 slices.
func CosineSimilarity(a, b []float32) float32 {
	dotProduct := math32.Dot(a, b)
	magnitudeA := Magnitude(a)
	magnitudeB := Magnitude(b)

	// Avoid division by zero
	if magnitudeA == 0 || magnitudeB == 0 {
		return 0
	}

	return dotProduct / (magnitudeA * magnitudeB)
}

// CosineDistance calculates the cosine distance between two float32 slices.
func CosineDistance(a, b []float32) float32 {
	return 1 - CosineSimilarity(a, b)
}
