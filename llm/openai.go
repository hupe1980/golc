package llm

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/util"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
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
	client *openai.Client
	opts   OpenAIOptions
}

func NewOpenAI(apiKey string, optFns ...func(o *OpenAIOptions)) (*OpenAI, error) {
	opts := OpenAIOptions{
		ModelName:        "text-davinci-002",
		Temperatur:       0.7,
		MaxTokens:        256,
		TopP:             1,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		N:                1,
		BatchSize:        20,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	openAI := &OpenAI{
		client: openai.NewClient(apiKey),
		opts:   opts,
	}

	return openAI, nil
}

func (o *OpenAI) GetTokenIDs(text string) ([]int, error) {
	e, err := tiktoken.EncodingForModel(o.opts.ModelName)
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

func (o *OpenAI) Generate(ctx context.Context, prompts []string) (*golc.LLMResult, error) {
	subPromps := util.ChunkBy(prompts, o.opts.BatchSize)

	choices := []openai.CompletionChoice{}
	tokenUsage := make(map[string]int)

	for _, prompt := range subPromps {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			res, err := o.client.CreateCompletion(ctx, openai.CompletionRequest{
				Prompt:      prompt,
				Model:       o.opts.ModelName,
				Temperature: o.opts.Temperatur,
				MaxTokens:   o.opts.MaxTokens,
				TopP:        o.opts.TopP,
			})
			if err != nil {
				return nil, err
			}

			choices = append(choices, res.Choices...)
			tokenUsage["completionTokens"] += res.Usage.CompletionTokens
			tokenUsage["promptTokens"] += res.Usage.PromptTokens
			tokenUsage["totalTokens"] += res.Usage.TotalTokens
		}
	}

	generations := util.Map(util.ChunkBy(choices, o.opts.N), func(promptChoices []openai.CompletionChoice, _ int) []golc.Generation {
		return util.Map(promptChoices, func(choice openai.CompletionChoice, _ int) golc.Generation {
			return golc.Generation{
				Text: choice.Text,
				Info: map[string]any{
					"FinishReason": choice.FinishReason,
					"Logprobs":     choice.LogProbs,
				},
			}
		})
	})

	return &golc.LLMResult{
		Generations: generations,
		LLMOutput: map[string]any{
			"modelName":  o.opts.ModelName,
			"tokenUsage": tokenUsage,
		},
	}, nil
}

func (o *OpenAI) GeneratePrompt(ctx context.Context, promptValues []golc.PromptValue) (*golc.LLMResult, error) {
	prompts := util.Map(promptValues, func(value golc.PromptValue, _ int) string {
		return value.String()
	})

	return o.Generate(ctx, prompts)
}

func (o *OpenAI) Predict(ctx context.Context, text string) (string, error) {
	result, err := o.Generate(ctx, []string{text})
	if err != nil {
		return "", err
	}

	return result.Generations[0][0].Text, nil
}

func (o *OpenAI) PredictMessages(ctx context.Context, messages []golc.ChatMessage) (golc.ChatMessage, error) {
	text, err := golc.StringifyChatMessages(messages)
	if err != nil {
		return nil, err
	}

	prediction, err := o.Predict(ctx, text)
	if err != nil {
		return nil, err
	}

	return golc.NewAIChatMessage(prediction), nil
}
