package math32

var (
	useAVX  bool // nolint unused
	useNEON bool // nolint unused
)

// Dot two vectors.
func Dot(a, b []float32) float32 {
	return dot(a, b)
}

func dotGeneric(a, b []float32) float32 {
	var ret float32
	for i := range a {
		ret += a[i] * b[i]
	}

	return ret
}

func SquaredL2(a, b []float32) float32 {
	return squaredL2(a, b)
}

func squaredL2Generic(a, b []float32) float32 {
	var distance float32
	for i := 0; i < len(a); i++ {
		distance += (a[i] - b[i]) * (a[i] - b[i])
	}

	return distance
}
