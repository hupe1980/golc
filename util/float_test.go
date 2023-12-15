package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat64ToFloat32(t *testing.T) {
	// Test case 1: Convert positive float64 values
	input1 := []float64{1.5, 2.8, 3.2}
	expected1 := []float32{1.5, 2.8, 3.2}
	assert.Equal(t, expected1, Float64ToFloat32(input1))

	// Test case 2: Convert negative float64 values
	input2 := []float64{-1.5, -2.8, -3.2}
	expected2 := []float32{-1.5, -2.8, -3.2}
	assert.Equal(t, expected2, Float64ToFloat32(input2))

	// Test case 3: Convert empty slice
	input3 := []float64{}
	expected3 := []float32{}
	assert.Equal(t, expected3, Float64ToFloat32(input3))
}

func TestFloat32ToFloat64(t *testing.T) {
	// Test case 1: Convert positive float32 values
	input1 := []float32{1.5, 2.8, 3.2}
	expected1 := []float64{1.5, 2.8, 3.2}
	assert.InEpsilonSlice(t, expected1, Float32ToFloat64(input1), 1e-07)

	// Test case 2: Convert negative float32 values
	input2 := []float32{-1.5, -2.8, -3.2}
	expected2 := []float64{-1.5, -2.8, -3.2}
	assert.InEpsilonSlice(t, expected2, Float32ToFloat64(input2), 1e-07)

	// Test case 3: Convert empty slice
	input3 := []float32{}
	expected3 := []float64{}
	assert.Equal(t, expected3, Float32ToFloat64(input3))
}
