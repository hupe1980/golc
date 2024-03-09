package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	t.Run("PutAndHas", func(t *testing.T) {
		set := NewSet[int]()
		set.Put(1)
		set.Put(2)
		set.Put(3)

		assert.True(t, set.Has(1))
		assert.True(t, set.Has(2))
		assert.True(t, set.Has(3))
		assert.False(t, set.Has(4))
	})

	t.Run("Remove", func(t *testing.T) {
		set := SetOf[int](1, 2, 3)
		set.Remove(2)

		assert.True(t, set.Has(1))
		assert.False(t, set.Has(2))
		assert.True(t, set.Has(3))
		assert.Equal(t, 2, set.Size())
	})

	t.Run("Clear", func(t *testing.T) {
		set := SetOf[int](1, 2, 3)
		set.Clear()

		assert.False(t, set.Has(1))
		assert.False(t, set.Has(2))
		assert.False(t, set.Has(3))
		assert.Equal(t, 0, set.Size())
	})

	t.Run("Each", func(t *testing.T) {
		set := SetOf[string]("apple", "banana", "cherry")

		var count int

		set.Each(func(key string) {
			count++
		})

		assert.Equal(t, 3, count)
	})

	t.Run("ToSlice", func(t *testing.T) {
		set := SetOf[int](1, 2, 3)
		slice := set.ToSlice()

		assert.ElementsMatch(t, []int{1, 2, 3}, slice)
	})
}
