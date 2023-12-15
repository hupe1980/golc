package memory

import (
	"context"

	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Simple satisfies the Memory interface.
var _ schema.Memory = (*Simple)(nil)

type Simple struct {
	memories map[string]any
}

func NewSimple() Simple {
	return Simple{
		memories: make(map[string]any),
	}
}

func (m *Simple) MemoryKeys() []string {
	return util.Keys(m.memories)
}

func (m *Simple) LoadMemoryVariables(ctx context.Context, inputs map[string]any) (map[string]any, error) {
	return m.memories, nil
}

func (m *Simple) SaveContext(ctx context.Context, inputs map[string]any, outputs map[string]any) error {
	return nil
}

func (m *Simple) Clear(ctx context.Context) error {
	return nil
}
