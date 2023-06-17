package memory

import "github.com/hupe1980/golc/schema"

// Compile time check to ensure Readonly satisfies the Memory interface.
var _ schema.Memory = (*Readonly)(nil)

type Readonly struct {
	memory schema.Memory
}

func NewReadonly(memory schema.Memory) Readonly {
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
