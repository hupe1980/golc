package memory

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/hupe1980/golc/chatmessagehistory"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure ConversationBuffer satisfies the Memory interface.
var _ schema.Memory = (*ConversationBuffer)(nil)

// ConversationBufferOptions contains options for configuring the ConversationBuffer memory type.
type ConversationBufferOptions struct {
	HumanPrefix        string
	AIPrefix           string
	MemoryKey          string
	InputKey           string
	OutputKey          string
	ReturnMessages     bool
	ChatMessageHistory schema.ChatMessageHistory

	// Size of the interactions window
	K uint
}

// ConversationBuffer is a memory type that manages conversation buffers.
type ConversationBuffer struct {
	opts ConversationBufferOptions
}

// NewConversationBuffer creates a new instance of ConversationBuffer memory type.
func NewConversationBuffer(optFns ...func(o *ConversationBufferOptions)) *ConversationBuffer {
	opts := ConversationBufferOptions{
		HumanPrefix:    "Human",
		AIPrefix:       "AI",
		MemoryKey:      "history",
		InputKey:       "",
		OutputKey:      "",
		ReturnMessages: false,
		K:              math.MaxUint,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.ChatMessageHistory == nil {
		opts.ChatMessageHistory = chatmessagehistory.NewInMemory()
	}

	return &ConversationBuffer{
		opts: opts,
	}
}

// MemoryKeys returns the memory keys for ConversationBuffer.
func (m *ConversationBuffer) MemoryKeys() []string {
	return []string{m.opts.MemoryKey}
}

// LoadMemoryVariables returns key-value pairs given the text input to the chain.
func (m *ConversationBuffer) LoadMemoryVariables(ctx context.Context, inputs map[string]any) (map[string]any, error) {
	messages, err := m.opts.ChatMessageHistory.Messages(ctx)
	if err != nil {
		return nil, err
	}

	if m.opts.K != math.MaxUint {
		if m.opts.K == 0 {
			messages = schema.ChatMessages{}
		} else {
			start := len(messages) - int(m.opts.K)*2
			if start > 0 {
				messages = messages[start:]
			}
		}
	}

	if m.opts.ReturnMessages {
		return map[string]any{
			m.opts.MemoryKey: messages,
		}, nil
	}

	buffer, err := messages.Format(func(o *schema.StringifyChatMessagesOptions) {
		o.HumanPrefix = m.opts.HumanPrefix
		o.AIPrefix = m.opts.AIPrefix
	})
	if err != nil {
		return nil, err
	}

	return map[string]any{
		m.opts.MemoryKey: buffer,
	}, nil
}

// SaveContext saves the input and output messages to the chat message history.
func (m *ConversationBuffer) SaveContext(ctx context.Context, inputs map[string]any, outputs map[string]any) error {
	input, output, err := m.getInputOutput(inputs, outputs)
	if err != nil {
		return err
	}

	if err := m.opts.ChatMessageHistory.AddUserMessage(ctx, input); err != nil {
		return err
	}

	return m.opts.ChatMessageHistory.AddAIMessage(ctx, output)
}

// Clear clears the chat message history.
func (m *ConversationBuffer) Clear(ctx context.Context) error {
	return m.opts.ChatMessageHistory.Clear(ctx)
}

func (m *ConversationBuffer) getInputOutput(inputs map[string]any, outputs map[string]any) (string, string, error) {
	inputKey := m.opts.InputKey
	if inputKey == "" {
		var err error

		inputKey, err = getPromptInputKey(inputs, m.MemoryKeys())
		if err != nil {
			return "", "", err
		}
	}

	input, ok := inputs[inputKey].(string)
	if !ok {
		return "", "", errors.New("")
	}

	outputKey := m.opts.OutputKey
	if outputKey == "" {
		if len(outputs) != 1 {
			return "", "", fmt.Errorf("multiple output keys. Only one output key expected, got %d", len(outputs))
		}

		for key := range outputs {
			outputKey = key
			break
		}
	}

	output, ok := outputs[outputKey].(string)
	if !ok {
		return "", "", errors.New("")
	}

	return input, output, nil
}

func getPromptInputKey(inputs map[string]interface{}, memoryVariables []string) (string, error) {
	promptInputKeys := make([]string, 0, len(inputs))

	for key := range inputs {
		if key != "stop" && !util.Contains(memoryVariables, key) {
			promptInputKeys = append(promptInputKeys, key)
		}
	}

	if len(promptInputKeys) != 1 {
		return "", fmt.Errorf("multiple input keys. One input key expected, got %d", len(promptInputKeys))
	}

	return promptInputKeys[0], nil
}
