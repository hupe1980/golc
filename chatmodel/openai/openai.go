package openai

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/chatmodel"
	"github.com/sashabaranov/go-openai"
)

type Options struct {
	// Model name to use.
	Model string
	// Sampling temperature to use.
	Temperatur float32
	// The maximum number of tokens to generate in the completion.
	// -1 returns as many tokens as possible given the prompt and
	//the models maximal context size.
	MaxTokens int
	// Total probability mass of tokens to consider at each step.
	TopP float32
	// Penalizes repeated tokens.
	PresencePenalty float32
	// Penalizes repeated tokens according to frequency.
	FrequencyPenalty float32
	// How many completions to generate for each prompt.
	N int
	// Batch size to use when passing multiple documents to generate.
	BatchSize int
}

type ChatOpenAI struct {
	*chatmodel.ChatModel
	client *openai.Client
	opts   Options
}

func New(apiKey string) (*ChatOpenAI, error) {
	opts := Options{
		Model:            "gpt-3.5-turbo",
		Temperatur:       1,
		TopP:             1,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
	}

	openai := &ChatOpenAI{
		client: openai.NewClient(apiKey),
		opts:   opts,
	}

	openai.ChatModel = chatmodel.NewChatModel(openai.generate)

	return openai, nil
}

func (o *ChatOpenAI) generate(ctx context.Context, messages []golc.ChatMessage) (*golc.LLMResult, error) {
	openAIMessages := []openai.ChatCompletionMessage{}

	for _, message := range messages {
		role, err := messageTypeToOpenAIRole(message.Type())
		if err != nil {
			return nil, err
		}

		openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
			Role:    role,
			Content: message.Text(),
		})
	}

	res, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    o.opts.Model,
		Messages: openAIMessages,
	})
	if err != nil {
		return nil, err
	}

	text := res.Choices[0].Message.Content
	role := res.Choices[0].Message.Role

	return &golc.LLMResult{
		Generations: [][]golc.Generation{{golc.Generation{
			Text:    text,
			Message: openAIResponseToChatMessage(role, text),
		}}},
		LLMOutput: map[string]any{},
	}, nil
}

func messageTypeToOpenAIRole(mType golc.ChatMessageType) (string, error) {
	switch mType { // nolint exhaustive
	case golc.ChatMessageTypeSystem:
		return "system", nil
	case golc.ChatMessageTypeAI:
		return "assistant", nil
	case golc.ChatMessageTypeHuman:
		return "user", nil
	default:
		return "", fmt.Errorf("unknown message type: %s", mType)
	}
}

func openAIResponseToChatMessage(role, text string) golc.ChatMessage {
	switch role {
	case "user":
		return golc.NewHumanChatMessage(text)
	case "assistant":
		return golc.NewAIChatMessage(text)
	case "system":
		return golc.NewSystemChatMessage(text)
	}

	return golc.NewGenericChatMessage(text, "unknown")
}
