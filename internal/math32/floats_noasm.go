//go:build noasm || (!arm64 && !amd64)

package math32

func dot(a, b []float32) float32 {
	return dotGeneric(a, b)
}

func squaredL2(a, b []float32) float32 {
	return squaredL2Generic(a, b)
}
