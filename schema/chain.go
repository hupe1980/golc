package schema

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hupe1980/golc/internal/util"
)

// ChainValues represents a map of key-value pairs used for passing inputs and outputs between chain components.
type ChainValues map[string]any

// GetString retrieves the value associated with the given name as a string from ChainValues.
// If the name does not exist or the value is not a string, an error is returned.
func (cv ChainValues) GetString(name string) (string, error) {
	value, ok := cv[name]
	if !ok {
		return "", fmt.Errorf("%w: no value for key %s", ErrInvalidChainValues, name)
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return "", ErrChainValueWrongType
	}
}

// GetDocuments retrieves the value associated with the given name as a slice of documents from ChainValues.
// If the name does not exist, the value is not a slice of documents, or the slice is empty, an error is returned.
func (cv ChainValues) GetDocuments(name string) ([]Document, error) {
	value, ok := cv[name]
	if !ok {
		return nil, fmt.Errorf("%w: no value for inputKey %s", ErrInvalidChainValues, name)
	}

	docs, ok := value.([]Document)
	if !ok {
		return nil, ErrChainValueWrongType
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("%w: no documents", ErrInvalidChainValues)
	}

	return docs, nil
}

// Clone creates a shallow copy of the ChainValues map.
func (cv ChainValues) Clone() ChainValues {
	return util.CopyMap(cv)
}

// CallOptions contains general options for executing a chain.
type CallOptions struct {
	CallbackManger CallbackManagerForChainRun
	Stop           []string
}

// Chain represents a sequence of calls to llms oder other utilities.
type Chain interface {
	// Call executes the chain with the given context and inputs.
	// It returns the outputs of the chain or an error, if any.
	Call(ctx context.Context, inputs ChainValues, optFns ...func(o *CallOptions)) (ChainValues, error)
	// Type returns the type of the chain.
	Type() string
	// Verbose returns the verbosity setting of the chain.
	Verbose() bool
	// Callbacks returns the callbacks associated with the chain.
	Callbacks() []Callback
	// Memory returns the memory associated with the chain.
	Memory() Memory
	// InputKeys returns the expected input keys.
	InputKeys() []string
	// OutputKeys returns the output keys the chain will return.
	OutputKeys() []string
}
