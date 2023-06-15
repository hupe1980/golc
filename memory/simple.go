package memory

import (
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/util"
)

// Compile time check to ensure Simple satisfies the memory interface.
var _ golc.Memory = (*Simple)(nil)

type Simple struct {
	memories map[string]any
}

func NewSimple() Simple {
	return Simple{
		memories: make(map[string]any),
	}
}

func (m *Simple) MemoryVariables() []string {
	return util.Keys(m.memories)
}

func (m *Simple) LoadMemoryVariables(inputs map[string]any) (map[string]any, error) {
	return m.memories, nil
}

func (m *Simple) SaveContext(inputs map[string]any, outputs map[string]any) error {
	return nil
}

func (m *Simple) Clear() error {
	return nil
}
