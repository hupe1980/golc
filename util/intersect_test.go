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
