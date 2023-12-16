package deepcopy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Field1 string
	Field2 int
}

func TestCopy(t *testing.T) {
	t.Run("Copy nil", func(t *testing.T) {
		result := Copy(nil)
		assert.Nil(t, result)
	})

	t.Run("Copy int", func(t *testing.T) {
		src := 42
		result := Copy(src)
		assert.Equal(t, src, result)
	})

	t.Run("Copy struct", func(t *testing.T) {
		src := TestStruct{Field1: "Test", Field2: 42}
		result := Copy(src)
		assert.Equal(t, src, result)
	})

	t.Run("Copy time.Time", func(t *testing.T) {
		src := time.Now()
		result := Copy(src)
		assert.Equal(t, src, result)
	})

	t.Run("Copy nested struct", func(t *testing.T) {
		src := struct {
			Nested TestStruct
		}{Nested: TestStruct{Field1: "Nested", Field2: 99}}

		result := Copy(src)
		assert.Equal(t, src, result)
	})

	t.Run("Copy slice", func(t *testing.T) {
		src := []int{1, 2, 3}
		result := Copy(src)
		assert.Equal(t, src, result)
	})

	t.Run("Copy map", func(t *testing.T) {
		src := map[string]int{"one": 1, "two": 2}
		result := Copy(src)
		assert.Equal(t, src, result)
	})
}
