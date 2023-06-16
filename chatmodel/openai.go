package chatmodel

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc"
	"github.com/sashabaranov/go-openai"

	"github.com/pkoukk/tiktoken-go"
)

// Compile time check to ensure OpenAI satisfies the llm interface.
var _ golc.LLM = (*OpenAI)(nil)

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

	openai := &OpenAI{
		client: openai.NewClient(apiKey),
		opts:   opts,
	}

	openai.ChatModel = NewChatModel(openai.generate)

	return openai, nil
}

func (o *OpenAI) generate(ctx context.Context, messages []golc.ChatMessage) (*golc.LLMResult, error) {
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

	return &golc.LLMResult{
		Generations: [][]golc.Generation{{golc.Generation{
			Text:    text,
			Message: openAIResponseToChatMessage(role, text),
		}}},
		LLMOutput: map[string]any{},
	}, nil
}

func (o *OpenAI) GetTokenIDs(text string) ([]int, error) {
	_, e, err := o.getEncodingForModel()
	if err != nil {
		return nil, err
	}

	return e.Encode(text, nil, nil), nil
}

func (o *OpenAI) GetNumTokens(text string) (int, error) {
	ids, err := o.GetTokenIDs(text)
	if err != nil {
		return 0, err
	}

	return len(ids), nil
}

func (o *OpenAI) GetNumTokensFromMessage(messages []golc.ChatMessage) (int, error) {
	text, err := golc.StringifyChatMessages(messages)
	if err != nil {
		return 0, err
	}

	return o.GetNumTokens(text)
}

func (o *OpenAI) getEncodingForModel() (string, *tiktoken.Tiktoken, error) {
	model := o.opts.ModelName
	if model == "gpt-3.5-turbo" {
		model = "gpt-3.5-turbo-0301"
	} else if model == "gpt-4" {
		model = "gpt-4-0314"
	}

	e, err := tiktoken.EncodingForModel(model)
	if err != nil {
		model = "cl100k_base" //fallback

		e, err = tiktoken.EncodingForModel(model)

		return model, e, err
	}

	return model, e, nil
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
