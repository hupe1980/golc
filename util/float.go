package util

// Float64ToFloat32 converts a slice of float64 values to a slice of float32 values.
// It creates a new slice and populates it with the corresponding float32 values.
func Float64ToFloat32(v []float64) []float32 {
	v32 := make([]float32, len(v))
	for i, f := range v {
		v32[i] = float32(f)
	}

	return v32
}

// Float32ToFloat64 converts a slice of float32 values to a slice of float64 values.
// It creates a new slice and populates it with the corresponding float64 values.
func Float32ToFloat64(v []float32) []float64 {
	v64 := make([]float64, len(v))
	for i, f := range v {
		v64[i] = float64(f)
	}

	return v64
}
