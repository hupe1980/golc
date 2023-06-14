package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOmitByKeys(t *testing.T) {
	m := OmitByKeys(map[string]int{"foo": 1, "bar": 2, "baz": 3}, []string{"foo", "baz"})
	assert.Equal(t, m, map[string]int{"bar": 2})
}

func TestKeys(t *testing.T) {
	keys := Keys(map[string]int{"foo": 1, "bar": 2})
	assert.ElementsMatch(t, keys, []string{"bar", "foo"})
}
