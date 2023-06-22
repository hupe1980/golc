package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	t.Run("ElementPresent", func(t *testing.T) {
		// Test case: Element is present in the collection
		collection := []int{1, 2, 3, 4, 5}
		element := 3

		result := Contains(collection, element)
		assert.True(t, result, "Element should be present in the collection")
	})

	t.Run("ElementNotPresent", func(t *testing.T) {
		// Test case: Element is not present in the collection
		collection := []string{"apple", "banana", "orange"}
		element := "grape"

		result := Contains(collection, element)
		assert.False(t, result, "Element should not be present in the collection")
	})

	t.Run("EmptyCollection", func(t *testing.T) {
		// Test case: Empty collection
		var collection []bool
		element := true

		result := Contains(collection, element)
		assert.False(t, result, "Element should not be present in an empty collection")
	})
}

func TestDifference(t *testing.T) {
	t.Run("IntSlice", func(t *testing.T) {
		list1 := []int{1, 2, 3, 4, 5}
		list2 := []int{4, 5, 6, 7, 8}

		left, right := Difference(list1, list2)
		expectedLeft := []int{1, 2, 3}
		expectedRight := []int{6, 7, 8}

		assert.ElementsMatch(t, expectedLeft, left)
		assert.ElementsMatch(t, expectedRight, right)
	})

	t.Run("StringSlice", func(t *testing.T) {
		list1 := []string{"apple", "banana", "cherry"}
		list2 := []string{"banana", "cherry", "date"}

		left, right := Difference(list1, list2)
		expectedLeft := []string{"apple"}
		expectedRight := []string{"date"}

		assert.ElementsMatch(t, expectedLeft, left)
		assert.ElementsMatch(t, expectedRight, right)
	})

	t.Run("EmptySlice", func(t *testing.T) {
		list1 := []int{}
		list2 := []int{1, 2, 3}

		left, right := Difference(list1, list2)
		expectedLeft := []int{}
		expectedRight := []int{1, 2, 3}

		assert.ElementsMatch(t, expectedLeft, left)
		assert.ElementsMatch(t, expectedRight, right)
	})

	t.Run("NoDifference", func(t *testing.T) {
		list1 := []int{1, 2, 3}
		list2 := []int{4, 5, 6}

		left, right := Difference(list1, list2)
		expectedLeft := []int{1, 2, 3}
		expectedRight := []int{4, 5, 6}

		assert.ElementsMatch(t, expectedLeft, left)
		assert.ElementsMatch(t, expectedRight, right)
	})
}

func TestIntersect(t *testing.T) {
	t.Run("IntSlice", func(t *testing.T) {
		list1 := []int{1, 2, 3, 4, 5}
		list2 := []int{4, 5, 6, 7, 8}

		intersection := Intersect(list1, list2)
		expected := []int{4, 5}

		assert.ElementsMatch(t, expected, intersection)
	})

	t.Run("StringSlice", func(t *testing.T) {
		list1 := []string{"apple", "banana", "cherry"}
		list2 := []string{"banana", "cherry", "date"}

		intersection := Intersect(list1, list2)
		expected := []string{"banana", "cherry"}

		assert.ElementsMatch(t, expected, intersection)
	})

	t.Run("EmptySlice", func(t *testing.T) {
		list1 := []int{}
		list2 := []int{1, 2, 3}

		intersection := Intersect(list1, list2)
		expected := []int{}

		assert.ElementsMatch(t, expected, intersection)
	})

	t.Run("NoIntersection", func(t *testing.T) {
		list1 := []int{1, 2, 3}
		list2 := []int{4, 5, 6}

		intersection := Intersect(list1, list2)
		expected := []int{}

		assert.ElementsMatch(t, expected, intersection)
	})
}
