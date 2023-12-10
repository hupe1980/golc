package memory

import (
	"context"
	"errors"
	"fmt"

	"github.com/hupe1980/golc/memory/chatmessagehistory"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure ConversationBuffer satisfies the Memory interface.
var _ schema.Memory = (*ConversationBuffer)(nil)

// ConversationTokenBufferOptions contains options for configuring the ConversationTokenBuffer memory type.
type ConversationTokenBufferOptions struct {
	HumanPrefix        string
	AIPrefix           string
	MemoryKey          string
	InputKey           string
	OutputKey          string
	ReturnMessages     bool
	ChatMessageHistory schema.ChatMessageHistory

	MaxTokenLimit uint
}

// ConversationTokenBuffer is a memory type that manages conversation token buffers.
type ConversationTokenBuffer struct {
	tokenizer schema.Tokenizer
	opts      ConversationTokenBufferOptions
}

// NewConversationTokenBuffer creates a new instance of ConversationTokenBuffer memory type.
func NewConversationTokenBuffer(tokenizer schema.Tokenizer, optFns ...func(o *ConversationTokenBufferOptions)) *ConversationTokenBuffer {
	opts := ConversationTokenBufferOptions{
		HumanPrefix:    "Human",
		AIPrefix:       "AI",
		MemoryKey:      "history",
		InputKey:       "",
		OutputKey:      "",
		ReturnMessages: false,
		MaxTokenLimit:  2000,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.ChatMessageHistory == nil {
		opts.ChatMessageHistory = chatmessagehistory.NewInMemory()
	}

	return &ConversationTokenBuffer{
		tokenizer: tokenizer,
		opts:      opts,
	}
}

// MemoryKeys returns the memory keys for ConversationTokenBuffer.
func (m *ConversationTokenBuffer) MemoryKeys() []string {
	return []string{m.opts.MemoryKey}
}

// LoadMemoryVariables returns key-value pairs given the text input to the chain.
func (m *ConversationTokenBuffer) LoadMemoryVariables(ctx context.Context, inputs map[string]any) (map[string]any, error) {
	messages, err := m.opts.ChatMessageHistory.Messages(ctx)
	if err != nil {
		return nil, err
	}

	numTokens, err := m.getNumTokensForMessages(messages)
	if err != nil {
		return nil, err
	}

	for len(messages) > 0 {
		if numTokens > m.opts.MaxTokenLimit {
			messages = messages[1:]

			numTokens, err = m.getNumTokensForMessages(messages)
			if err != nil {
				return nil, err
			}
		} else {
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
	}

	if m.opts.ReturnMessages {
		return map[string]any{
			m.opts.MemoryKey: schema.ChatMessages{},
		}, nil
	}

	return map[string]any{
		m.opts.MemoryKey: "",
	}, nil
}

// SaveContext saves the input and output messages to the chat message history.
func (m *ConversationTokenBuffer) SaveContext(ctx context.Context, inputs map[string]any, outputs map[string]any) error {
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
func (m *ConversationTokenBuffer) Clear(ctx context.Context) error {
	return m.opts.ChatMessageHistory.Clear(ctx)
}

func (m *ConversationTokenBuffer) getInputOutput(inputs map[string]any, outputs map[string]any) (string, string, error) {
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

func (m *ConversationTokenBuffer) getNumTokensForMessages(messages schema.ChatMessages) (uint, error) {
	buffer, err := messages.Format(func(o *schema.StringifyChatMessagesOptions) {
		o.HumanPrefix = m.opts.HumanPrefix
		o.AIPrefix = m.opts.AIPrefix
	})
	if err != nil {
		return 0, err
	}

	return m.tokenizer.GetNumTokens(buffer)
}
