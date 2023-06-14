package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	res := Map([]int{1, 2, 3, 4}, func(x int, _ int) string {
		return "xxx"
	})

	assert.Equal(t, res, []string{"xxx", "xxx", "xxx", "xxx"})
}

func TestFilter(t *testing.T) {
	f := Filter([]string{"", "foo", "", "bar", ""}, func(x string, _ int) bool {
		return len(x) > 2
	})
	assert.ElementsMatch(t, f, []string{"foo", "bar"})
}
