package llm

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
	"github.com/hupe1980/golc/util"
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
	callbackOptions
}

type OpenAI struct {
	*llm
	schema.Tokenizer
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
		callbackOptions: callbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	openAI := &OpenAI{
		Tokenizer: tokenizer.NewOpenAI(opts.ModelName),
		client:    openai.NewClient(apiKey),
		opts:      opts,
	}

	openAI.llm = newLLM("OpenAI", openAI.generate, opts.Verbose)

	return openAI, nil
}

func (o *OpenAI) generate(ctx context.Context, prompts []string, stop []string) (*schema.LLMResult, error) {
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
				Stop:        stop,
			})
			if err != nil {
				return nil, err
			}

			choices = append(choices, res.Choices...)
			tokenUsage["CompletionTokens"] += res.Usage.CompletionTokens
			tokenUsage["PromptTokens"] += res.Usage.PromptTokens
			tokenUsage["TotalTokens"] += res.Usage.TotalTokens
		}
	}

	generations := util.Map(util.ChunkBy(choices, o.opts.N), func(promptChoices []openai.CompletionChoice, _ int) []*schema.Generation {
		return util.Map(promptChoices, func(choice openai.CompletionChoice, _ int) *schema.Generation {
			return &schema.Generation{
				Text: choice.Text,
				Info: map[string]any{
					"FinishReason": choice.FinishReason,
					"Logprobs":     choice.LogProbs,
				},
			}
		})
	})

	return &schema.LLMResult{
		Generations: generations,
		LLMOutput: map[string]any{
			"ModelName":  o.opts.ModelName,
			"TokenUsage": tokenUsage,
		},
	}, nil
}
