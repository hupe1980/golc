package memory

import (
	"github.com/hupe1980/golc"
)

// Compile time check to ensure Readonly satisfies the memory interface.
var _ golc.Memory = (*Readonly)(nil)

type Readonly struct {
	memory golc.Memory
}

func NewReadonly(memory golc.Memory) Readonly {
	return Readonly{
		memory: memory,
	}
}

func (m *Readonly) MemoryVariables() []string {
	return m.memory.MemoryVariables()
}

func (m *Readonly) LoadMemoryVariables(inputs map[string]any) (map[string]any, error) {
	return m.memory.LoadMemoryVariables(inputs)
}

func (m *Readonly) SaveContext(inputs map[string]any, outputs map[string]any) error {
	return nil
}

func (m *Readonly) Clear() error {
	return nil
}
