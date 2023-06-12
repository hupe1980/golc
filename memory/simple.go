package memory

import (
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/util"
)

type Simple struct {
	memories map[string]any
}

func NewSimple() Simple {
	return Simple{
		memories: make(map[string]any),
	}
}

// Compile time check to ensure Simple satisfies the memory interface.
var _ golc.Memory = (*Simple)(nil)

func (m *Simple) MemoryVariables() []string {
	return util.Keys(m.memories)
}

func (m *Simple) LoadMemoryVariables(map[string]any) (map[string]any, error) {
	return m.memories, nil
}

func (m *Simple) SaveContext(map[string]any, map[string]any) error {
	return nil
}

func (m *Simple) Clear() error {
	return nil
}
