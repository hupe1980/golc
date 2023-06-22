package memory

import (
	"context"

	"github.com/hupe1980/golc/schema"
)

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

func (m *Readonly) MemoryKeys() []string {
	return m.memory.MemoryKeys()
}

func (m *Readonly) LoadMemoryVariables(ctx context.Context, inputs map[string]any) (map[string]any, error) {
	return m.memory.LoadMemoryVariables(ctx, inputs)
}

func (m *Readonly) SaveContext(ctx context.Context, inputs map[string]any, outputs map[string]any) error {
	return nil
}

func (m *Readonly) Clear(ctx context.Context) error {
	return nil
}
