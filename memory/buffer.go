package memory

import (
	"context"
	"errors"
	"fmt"

	"github.com/hupe1980/golc/memory/chatmessagehistory"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
)

// Compile time check to ensure ConversationBuffer satisfies the Memory interface.
var _ schema.Memory = (*ConversationBuffer)(nil)

type ConversationBufferOptions struct {
	HumanPrefix        string
	AIPrefix           string
	MemoryKey          string
	InputKey           string
	OutputKey          string
	ReturnMessages     bool
	ChatMessageHistory schema.ChatMessageHistory
}

type ConversationBuffer struct {
	opts ConversationBufferOptions
}

func NewConversationBuffer(optFns ...func(o *ConversationBufferOptions)) *ConversationBuffer {
	opts := ConversationBufferOptions{
		HumanPrefix:    "Human",
		AIPrefix:       "AI",
		MemoryKey:      "history",
		InputKey:       "",
		OutputKey:      "",
		ReturnMessages: false,
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

func (m *ConversationBuffer) MemoryKeys() []string {
	return []string{m.opts.MemoryKey}
}

func (m *ConversationBuffer) LoadMemoryVariables(ctx context.Context, inputs map[string]any) (map[string]any, error) {
	messages, err := m.opts.ChatMessageHistory.Messages(ctx)
	if err != nil {
		return nil, err
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

func (m *ConversationBuffer) SaveContext(ctx context.Context, inputs map[string]any, outputs map[string]any) error {
	input, output, err := m.getInputOutput(inputs, outputs)
	if err != nil {
		return err
	}

	if err := m.opts.ChatMessageHistory.AddUserMessage(ctx, input); err != nil {
		return err
	}

	if err := m.opts.ChatMessageHistory.AddAIMessage(ctx, output); err != nil {
		return err
	}

	return nil
}

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
