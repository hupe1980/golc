package integration

import (
	"fmt"

	"github.com/hupe1980/golc/schema"
	"github.com/sashabaranov/go-openai"
)

// ToOpenAIChatCompletionMessages converts a slice of schema.ChatMessage to a slice of openai.ChatCompletionMessage.
// It extracts the necessary information from each message to create the corresponding OpenAI chat completion message.
func ToOpenAIChatCompletionMessages(messages schema.ChatMessages) ([]openai.ChatCompletionMessage, error) {
	openAIMessages := []openai.ChatCompletionMessage{}

	for _, message := range messages {
		role, err := messageTypeToOpenAIRole(message.Type())
		if err != nil {
			return nil, err
		}

		if functionMessage, ok := message.(*schema.FunctionChatMessage); ok {
			openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
				Role:    role,
				Content: functionMessage.Content(),
				Name:    functionMessage.Name(),
			})
		} else {
			openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
				Role:    role,
				Content: message.Content(),
			})
		}
	}

	return openAIMessages, nil
}

// messageTypeToOpenAIRole converts a schema.ChatMessageType to the corresponding OpenAI role string.
func messageTypeToOpenAIRole(mType schema.ChatMessageType) (string, error) {
	switch mType { // nolint exhaustive
	case schema.ChatMessageTypeSystem:
		return "system", nil
	case schema.ChatMessageTypeAI:
		return "assistant", nil
	case schema.ChatMessageTypeHuman:
		return "user", nil
	case schema.ChatMessageTypeFunction:
		return "function", nil
	default:
		return "", fmt.Errorf("unknown message type: %s", mType)
	}
}
