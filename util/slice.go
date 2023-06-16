package util

// Map manipulates a slice and transforms it to a slice of another type.
func Map[T, U any](ts []T, f func(e T, i int) U) []U {
	res := make([]U, len(ts))
	for i, e := range ts {
		res[i] = f(e, i)
	}

	return res
}

// ChunkBy splits a slice into chunks of a specified size.
func ChunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}

	return append(chunks, items)
}

// Filter applies a filtering function to a collection and returns a new slice
// containing the elements that satisfy the provided predicate function.
func Filter[T any](collection []T, f func(e T, i int) bool) []T {
	fltd := make([]T, 0, len(collection))

	for i, e := range collection {
		if f(e, i) {
			fltd = append(fltd, e)
		}
	}

	return fltd
}

// SumInt calculates the sum of all integers in the given slice.
func SumInt(slice []int) int {
	total := 0
	for _, num := range slice {
		total += num
	}

	return total
}
