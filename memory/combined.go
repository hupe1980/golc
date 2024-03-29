package memory

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Combined satisfies the Memory interface.
var _ schema.Memory = (*Combined)(nil)

type Combined struct {
	memories []schema.Memory
}

func NewCombined(memories ...schema.Memory) (*Combined, error) {
	if err := checkRepeatedMemoryVariable(memories...); err != nil {
		return nil, err
	}

	return &Combined{
		memories: memories,
	}, nil
}

func (m *Combined) MemoryKeys() []string {
	memoryKeys := make([]string, 0)
	for _, memory := range m.memories {
		memoryKeys = append(memoryKeys, memory.MemoryKeys()...)
	}

	return memoryKeys
}

func (m *Combined) LoadMemoryVariables(ctx context.Context, inputs map[string]any) (map[string]any, error) {
	memoryData := make(map[string]any)

	// Collect vars from the sub-memories
	for _, memory := range m.memories {
		data, err := memory.LoadMemoryVariables(ctx, inputs)
		if err != nil {
			return nil, err
		}

		for key, value := range data {
			memoryData[key] = value
		}
	}

	return memoryData, nil
}

func (m *Combined) SaveContext(ctx context.Context, inputs map[string]any, outputs map[string]any) error {
	for _, memory := range m.memories {
		if err := memory.SaveContext(ctx, inputs, outputs); err != nil {
			return err
		}
	}

	return nil
}

func (m *Combined) Clear(ctx context.Context) error {
	for _, memory := range m.memories {
		if err := memory.Clear(ctx); err != nil {
			return err
		}
	}

	return nil
}

func checkRepeatedMemoryVariable(memories ...schema.Memory) error {
	allKeys := make(map[string]bool)

	for _, m := range memories {
		for _, key := range m.MemoryKeys() {
			if allKeys[key] {
				return fmt.Errorf("repeated memory key found: %s", key)
			}

			allKeys[key] = true
		}
	}

	return nil
}
