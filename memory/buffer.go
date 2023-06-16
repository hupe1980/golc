package memory

import (
	"errors"
	"fmt"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/memory/chatmessagehistory"
	"github.com/hupe1980/golc/util"
)

// Compile time check to ensure ConversationBuffer satisfies the memory interface.
var _ golc.Memory = (*ConversationBuffer)(nil)

type ConversationBufferOptions struct {
	HumanPrefix        string
	AIPrefix           string
	MemoryKey          string
	InputKey           string
	OutputKey          string
	ReturnMessages     bool
	ChatMessageHistory golc.ChatMessageHistory
}

type ConversationBuffer struct {
	opts ConversationBufferOptions
}

func NewConversationBuffer() *ConversationBuffer {
	opts := ConversationBufferOptions{
		HumanPrefix:    "Human",
		AIPrefix:       "AI",
		MemoryKey:      "history",
		InputKey:       "",
		OutputKey:      "",
		ReturnMessages: false,
	}

	if opts.ChatMessageHistory == nil {
		opts.ChatMessageHistory = chatmessagehistory.NewInMemory()
	}

	return &ConversationBuffer{
		opts: opts,
	}
}

func (m *ConversationBuffer) MemoryVariables() []string {
	return []string{m.opts.MemoryKey}
}

func (m *ConversationBuffer) LoadMemoryVariables(inputs map[string]any) (map[string]any, error) {
	messages, err := m.opts.ChatMessageHistory.Messages()
	if err != nil {
		return nil, err
	}

	if m.opts.ReturnMessages {
		return map[string]any{
			m.opts.MemoryKey: messages,
		}, nil
	}

	buffer, err := golc.StringifyChatMessages(messages, func(o *golc.StringifyChatMessagesOptions) {
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

func (m *ConversationBuffer) SaveContext(inputs map[string]any, outputs map[string]any) error {
	input, output, err := m.getInputOutput(inputs, outputs)
	if err != nil {
		return err
	}

	if err := m.opts.ChatMessageHistory.AddUserMessage(input); err != nil {
		return err
	}

	if err := m.opts.ChatMessageHistory.AddAIMessage(output); err != nil {
		return err
	}

	return nil
}

func (m *ConversationBuffer) Clear() error {
	return m.opts.ChatMessageHistory.Clear()
}

func (m *ConversationBuffer) getInputOutput(inputs map[string]any, outputs map[string]any) (string, string, error) {
	inputKey := m.opts.InputKey
	if inputKey == "" {
		var err error

		inputKey, err = getPromptInputKey(inputs, m.MemoryVariables())
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
