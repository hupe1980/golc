package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeMaps(t *testing.T) {
	t.Run("MergeTwoMaps", func(t *testing.T) {
		// Test case: Merge two maps into a single map
		map1 := map[string]int{"a": 1, "b": 2}
		map2 := map[string]int{"c": 3, "d": 4}

		result := MergeMaps(map1, map2)

		expectedResult := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

		assert.Equal(t, expectedResult, result, "Merged map is not as expected")
	})

	t.Run("MergeThreeMaps", func(t *testing.T) {
		// Test case: Merge three maps into a single map
		map1 := map[string]string{"a": "apple"}
		map2 := map[string]string{"b": "banana"}
		map3 := map[string]string{"c": "cherry"}

		result := MergeMaps(map1, map2, map3)

		expectedResult := map[string]string{"a": "apple", "b": "banana", "c": "cherry"}

		assert.Equal(t, expectedResult, result, "Merged map is not as expected")
	})

	t.Run("MergeEmptyMaps", func(t *testing.T) {
		// Test case: Merge empty maps
		var (
			map1 map[string]int
			map2 map[string]int
		)

		result := MergeMaps(map1, map2)

		expectedResult := map[string]int{}

		assert.Equal(t, expectedResult, result, "Merged map is not as expected")
	})
}

func TestCopyMap(t *testing.T) {
	t.Run("CopyNonEmptyMap", func(t *testing.T) {
		// Test case: Copy a non-empty map
		map1 := map[string]int{"a": 1, "b": 2, "c": 3}

		result := CopyMap(map1)

		assert.Equal(t, map1, result, "Copied map is not equal to the original map")
		assert.False(t, &map1 == &result, "Copied map should be distinct from the original map")
	})

	t.Run("CopyEmptyMap", func(t *testing.T) {
		// Test case: Copy an empty map
		var map1 map[string]int

		result := CopyMap(map1)

		expectedResult := map[string]int{}

		assert.Equal(t, expectedResult, result, "Copied map is not equal to an empty map")
	})
}

func TestOmitByKeys(t *testing.T) {
	t.Run("OmitExistingKeys", func(t *testing.T) {
		// Test case: Omit existing keys from the map
		map1 := map[string]int{"a": 1, "b": 2, "c": 3}
		keys := []string{"a", "c"}

		result := OmitByKeys(map1, keys)

		expectedResult := map[string]int{"b": 2}

		assert.Equal(t, expectedResult, result, "Omitted map is not as expected")
	})

	t.Run("OmitNonExistingKeys", func(t *testing.T) {
		// Test case: Omit non-existing keys from the map
		map1 := map[string]int{"a": 1, "b": 2, "c": 3}
		keys := []string{"d", "e"}

		result := OmitByKeys(map1, keys)

		expectedResult := map[string]int{"a": 1, "b": 2, "c": 3}

		assert.Equal(t, expectedResult, result, "Omitted map is not as expected")
	})

	t.Run("OmitAllKeys", func(t *testing.T) {
		// Test case: Omit all keys from the map
		map1 := map[string]int{"a": 1, "b": 2, "c": 3}
		keys := []string{"a", "b", "c"}

		result := OmitByKeys(map1, keys)

		expectedResult := map[string]int{}

		assert.Equal(t, expectedResult, result, "Omitted map is not as expected")
	})

	t.Run("OmitEmptyMap", func(t *testing.T) {
		// Test case: Omit keys from an empty map
		var map1 map[string]int

		keys := []string{"a", "b"}

		result := OmitByKeys(map1, keys)

		expectedResult := map[string]int{}

		assert.Equal(t, expectedResult, result, "Omitted map is not as expected")
	})
}

func TestKeys(t *testing.T) {
	t.Run("GetKeysFromMap", func(t *testing.T) {
		// Test case: Get keys from a map
		map1 := map[string]int{"a": 1, "b": 2, "c": 3}

		result := Keys(map1)

		expectedResult := []string{"a", "b", "c"}

		assert.ElementsMatch(t, expectedResult, result, "Keys are not as expected")
	})

	t.Run("GetKeysFromEmptyMap", func(t *testing.T) {
		// Test case: Get keys from an empty map
		var map1 map[string]int

		result := Keys(map1)

		expectedResult := []string{}

		assert.ElementsMatch(t, expectedResult, result, "Keys are not as expected")
	})
}

func TestKeyDifference(t *testing.T) {
	// Test case 1: Map1 and Map2 are empty
	t.Run("Empty maps", func(t *testing.T) {
		map1 := map[string]interface{}{}
		map2 := map[string]interface{}{}
		expected := []string{}
		difference := KeyDifference(map1, map2)
		assert.ElementsMatch(t, expected, difference, "Unexpected difference in keys")
	})

	// Test case 2: Map1 is empty, Map2 has keys
	t.Run("Map1 is empty", func(t *testing.T) {
		map1 := map[string]interface{}{}
		map2 := map[string]interface{}{
			"key1": 1,
			"key2": "value2",
		}
		expected := []string{"key1", "key2"}
		difference := KeyDifference(map1, map2)
		assert.ElementsMatch(t, expected, difference, "Unexpected difference in keys")
	})

	// Test case 3: Map1 has keys, Map2 is empty
	t.Run("Map2 is empty", func(t *testing.T) {
		map1 := map[string]interface{}{
			"key1": 1,
			"key2": "value2",
		}
		map2 := map[string]interface{}{}
		expected := []string{"key1", "key2"}
		difference := KeyDifference(map1, map2)
		assert.ElementsMatch(t, expected, difference, "Unexpected difference in keys")
	})

	// Test case 4: Map1 and Map2 have the same keys
	t.Run("Same keys", func(t *testing.T) {
		map1 := map[string]interface{}{
			"key1": 1,
			"key2": "value2",
		}
		map2 := map[string]interface{}{
			"key1": 2,
			"key2": "value2",
		}
		expected := []string{}
		difference := KeyDifference(map1, map2)
		assert.ElementsMatch(t, expected, difference, "Unexpected difference in keys")
	})

	// Test case 5: Map1 and Map2 have different keys
	t.Run("Different keys", func(t *testing.T) {
		map1 := map[string]interface{}{
			"key1": 1,
			"key2": "value2",
		}
		map2 := map[string]interface{}{
			"key3": true,
			"key4": 4.5,
		}
		expected := []string{"key1", "key2", "key3", "key4"}
		difference := KeyDifference(map1, map2)
		assert.ElementsMatch(t, expected, difference, "Unexpected difference in keys")
	})

	t.Run("Default", func(t *testing.T) {
		map1 := map[string]interface{}{
			"key1": 1,
			"key2": "value2",
			"key3": true,
		}

		map2 := map[string]interface{}{
			"key1": 1,
			"key3": true,
			"key4": 4.5,
		}

		expected := []string{"key2", "key4"}

		difference := KeyDifference(map1, map2)
		assert.ElementsMatch(t, expected, difference, "Unexpected difference in keys")
	})
}
