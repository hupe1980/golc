package math32

import "math"

// Sqrt returns the square root of x.
func Sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}
