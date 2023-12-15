package util

// MergeMaps merges multiple maps into a single map and returns the merged result.
func MergeMaps[M ~map[K]V, K comparable, V any](src ...M) M {
	merged := make(M)

	for _, m := range src {
		for k, v := range m {
			merged[k] = v
		}
	}

	return merged
}

// CopyMap creates a new copy of the given map and returns it.
func CopyMap[K, V comparable](m map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		result[k] = v
	}

	return result
}

// OmitByKeys creates a new map by omitting key-value pairs from
// the input map based on the specified keys.
func OmitByKeys[K comparable, V any](in map[K]V, keys []K) map[K]V {
	r := map[K]V{}

	for k, v := range in {
		if !Contains(keys, k) {
			r[k] = v
		}
	}

	return r
}

// Keys returns the keys of the map m.
// The keys will be an indeterminate order.
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}

	return r
}

// KeyDifference finds the keys that are present in one map
// and not in the other map.
func KeyDifference(map1, map2 map[string]any) []string {
	keys1 := make(map[string]bool)
	keys2 := make(map[string]bool)

	for key := range map1 {
		keys1[key] = true
	}

	for key := range map2 {
		keys2[key] = true
	}

	difference := make([]string, 0)

	for key := range keys1 {
		if !keys2[key] {
			difference = append(difference, key)
		}
	}

	for key := range keys2 {
		if !keys1[key] {
			difference = append(difference, key)
		}
	}

	return difference
}
