package util

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	t.Run("MapStringToInt", func(t *testing.T) {
		// Test case: Map string slice to int slice
		strSlice := []string{"1", "2", "3"}

		result := Map(strSlice, func(e string, i int) int {
			return len(e)
		})

		expectedResult := []int{1, 1, 1}

		assert.ElementsMatch(t, expectedResult, result, "Mapped slice is not as expected")
	})

	t.Run("MapEmptySlice", func(t *testing.T) {
		// Test case: Map an empty slice
		var emptySlice []int

		result := Map(emptySlice, func(e int, i int) bool {
			return e > 0
		})

		expectedResult := []bool{}

		assert.ElementsMatch(t, expectedResult, result, "Mapped slice is not as expected")
	})

	t.Run("MapWithIndex", func(t *testing.T) {
		// Test case: Map slice with access to the index
		intSlice := []int{1, 2, 3}

		result := Map(intSlice, func(e int, i int) string {
			return fmt.Sprintf("%d", e+i)
		})

		expectedResult := []string{"1", "3", "5"}

		assert.ElementsMatch(t, expectedResult, result, "Mapped slice is not as expected")
	})
}

func TestChunkBy(t *testing.T) {
	t.Run("ChunkByValidSize", func(t *testing.T) {
		// Test case: Chunk with a valid chunk size
		items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		chunkSize := 3

		result := ChunkBy(items, chunkSize)

		expectedResult := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}

		assert.Equal(t, expectedResult, result, "Chunks are not as expected")
	})

	t.Run("ChunkBySmallerSize", func(t *testing.T) {
		// Test case: Chunk with a size smaller than the number of items
		items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		chunkSize := 5

		result := ChunkBy(items, chunkSize)

		expectedResult := [][]int{{1, 2, 3, 4, 5}, {6, 7, 8, 9}}

		assert.Equal(t, expectedResult, result, "Chunks are not as expected")
	})

	t.Run("ChunkByLargerSize", func(t *testing.T) {
		// Test case: Chunk with a size larger than the number of items
		items := []int{1, 2, 3, 4, 5}
		chunkSize := 10

		result := ChunkBy(items, chunkSize)

		expectedResult := [][]int{{1, 2, 3, 4, 5}}

		assert.Equal(t, expectedResult, result, "Chunks are not as expected")
	})

	t.Run("ChunkByEmptySlice", func(t *testing.T) {
		// Test case: Chunk an empty slice
		var emptySlice []int
		chunkSize := 3

		result := ChunkBy(emptySlice, chunkSize)

		expectedResult := [][]int{[]int(nil)}

		assert.Equal(t, expectedResult, result, "Chunks are not as expected")
	})
}

func TestFilter(t *testing.T) {
	t.Run("FilterEvenNumbers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		evenNumbers := Filter(numbers, func(e int, i int) bool {
			return e%2 == 0
		})
		expectedEvenNumbers := []int{2, 4, 6, 8}
		assert.Equal(t, expectedEvenNumbers, evenNumbers, "Filtered even numbers should match expected result")
	})

	t.Run("FilterWordsStartingWithA", func(t *testing.T) {
		words := []string{"apple", "banana", "cherry", "date"}
		aWords := Filter(words, func(e string, i int) bool {
			return strings.HasPrefix(e, "a")
		})
		expectedAWords := []string{"apple"}
		assert.Equal(t, expectedAWords, aWords, "Filtered words starting with 'a' should match expected result")
	})

	t.Run("FilterPositiveNumbers", func(t *testing.T) {
		numbers2 := []int{-2, -1, 0, 1, 2}
		positiveNumbers := Filter(numbers2, func(e int, i int) bool {
			return e > 0
		})
		expectedPositiveNumbers := []int{1, 2}
		assert.Equal(t, expectedPositiveNumbers, positiveNumbers, "Filtered positive numbers should match expected result")
	})
}

func TestSumInt(t *testing.T) {
	t.Run("Empty Slice", func(t *testing.T) {
		slice := []int{}
		expectedSum := 0

		result := SumInt(slice)

		assert.Equal(t, expectedSum, result, "Sum of an empty slice should be 0")
	})

	t.Run("Positive Integers", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		expectedSum := 15

		result := SumInt(slice)

		assert.Equal(t, expectedSum, result, "Sum of positive integers should be calculated correctly")
	})

	t.Run("Negative Integers", func(t *testing.T) {
		slice := []int{-1, -2, -3, -4, -5}
		expectedSum := -15

		result := SumInt(slice)

		assert.Equal(t, expectedSum, result, "Sum of negative integers should be calculated correctly")
	})

	t.Run("Mixed Integers", func(t *testing.T) {
		slice := []int{-5, 2, -8, 10, -3}
		expectedSum := -4

		result := SumInt(slice)

		assert.Equal(t, expectedSum, result, "Sum of mixed positive and negative integers should be calculated correctly")
	})
}
