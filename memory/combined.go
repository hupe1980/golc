package memory

import (
	"fmt"

	"github.com/hupe1980/golc"
)

// Compile time check to ensure Combined satisfies the memory interface.
var _ golc.Memory = (*Combined)(nil)

type Combined struct {
	memories []golc.Memory
}

func NewCombined(memories ...golc.Memory) (*Combined, error) {
	if err := checkRepeatedMemoryVariable(memories...); err != nil {
		return nil, err
	}

	return &Combined{
		memories: memories,
	}, nil
}

func (m *Combined) MemoryVariables() []string {
	memoryVariables := make([]string, 0)
	for _, memory := range m.memories {
		memoryVariables = append(memoryVariables, memory.MemoryVariables()...)
	}

	return memoryVariables
}

func (m *Combined) LoadMemoryVariables(inputs map[string]any) (map[string]any, error) {
	memoryData := make(map[string]any)

	// Collect vars from the sub-memories
	for _, memory := range m.memories {
		data, err := memory.LoadMemoryVariables(inputs)
		if err != nil {
			return nil, err
		}

		for key, value := range data {
			memoryData[key] = value
		}
	}

	return memoryData, nil
}

func (m *Combined) SaveContext(inputs map[string]any, outputs map[string]any) error {
	for _, memory := range m.memories {
		if err := memory.SaveContext(inputs, outputs); err != nil {
			return err
		}
	}

	return nil
}

func (m *Combined) Clear() error {
	for _, memory := range m.memories {
		if err := memory.Clear(); err != nil {
			return err
		}
	}

	return nil
}

func checkRepeatedMemoryVariable(memories ...golc.Memory) error {
	allVariables := make(map[string]bool)

	for _, m := range memories {
		for _, variable := range m.MemoryVariables() {
			if allVariables[variable] {
				return fmt.Errorf("repeated memory variable found: %s", variable)
			}

			allVariables[variable] = true
		}
	}

	return nil
}
