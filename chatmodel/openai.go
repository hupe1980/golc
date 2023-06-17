package chatmodel

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
	"github.com/sashabaranov/go-openai"
)

// Compile time check to ensure OpenAI satisfies the LLM interface.
var _ schema.LLM = (*OpenAI)(nil)

type OpenAIOptions struct {
	// Model name to use.
	ModelName string
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

type OpenAI struct {
	*ChatModel
	schema.Tokenizer
	client *openai.Client
	opts   OpenAIOptions
}

func NewOpenAI(apiKey string) (*OpenAI, error) {
	opts := OpenAIOptions{
		ModelName:        "gpt-3.5-turbo",
		Temperatur:       1,
		TopP:             1,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
	}

	o := &OpenAI{
		Tokenizer: tokenizer.NewOpenAI(opts.ModelName),
		client:    openai.NewClient(apiKey),
		opts:      opts,
	}

	o.ChatModel = NewChatModel(o.generate)

	return o, nil
}

func (o *OpenAI) generate(ctx context.Context, messages []schema.ChatMessage, optFns ...func(o *schema.GenerateOptions)) (*schema.LLMResult, error) {
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
		Model:    o.opts.ModelName,
		Messages: openAIMessages,
	})
	if err != nil {
		return nil, err
	}

	text := res.Choices[0].Message.Content
	role := res.Choices[0].Message.Role

	return &schema.LLMResult{
		Generations: [][]*schema.Generation{{&schema.Generation{
			Text:    text,
			Message: openAIResponseToChatMessage(role, text),
		}}},
		LLMOutput: map[string]any{},
	}, nil
}

func messageTypeToOpenAIRole(mType schema.ChatMessageType) (string, error) {
	switch mType { // nolint exhaustive
	case schema.ChatMessageTypeSystem:
		return "system", nil
	case schema.ChatMessageTypeAI:
		return "assistant", nil
	case schema.ChatMessageTypeHuman:
		return "user", nil
	default:
		return "", fmt.Errorf("unknown message type: %s", mType)
	}
}

func openAIResponseToChatMessage(role, text string) schema.ChatMessage {
	switch role {
	case "user":
		return schema.NewHumanChatMessage(text)
	case "assistant":
		return schema.NewAIChatMessage(text)
	case "system":
		return schema.NewSystemChatMessage(text)
	}

	return schema.NewGenericChatMessage(text, "unknown")
}
