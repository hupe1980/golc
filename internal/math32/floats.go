package math32

var (
	useAVX512 bool // nolint unused
	useNEON   bool // nolint unused
)

// Dot two vectors.
func Dot(a, b []float32) (ret float32) {
	if len(a) != len(b) {
		panic("slice lengths do not match")
	}

	return dot(a, b)
}

func dotGeneric(a, b []float32) float32 {
	var ret float32
	for i := range a {
		ret += a[i] * b[i]
	}

	return ret
}
