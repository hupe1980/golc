package llm

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/util"
	openai "github.com/sashabaranov/go-openai"
)

// Compile time check to ensure OpenAI satisfies the llm interface.
var _ golc.LLM = (*OpenAI)(nil)

type OpenAIOptions struct {
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

type OpenAI struct {
	*LLM
	client *openai.Client
	opts   OpenAIOptions
}

func NewOpenAI(apiKey string) (*OpenAI, error) {
	opts := OpenAIOptions{
		Model:            "text-davinci-002",
		Temperatur:       0.7,
		MaxTokens:        256,
		TopP:             1,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		N:                1,
		BatchSize:        20,
	}

	openAI := &OpenAI{
		client: openai.NewClient(apiKey),
		opts:   opts,
	}

	openAI.LLM = NewLLM(openAI.generate)

	return openAI, nil
}

func (o *OpenAI) generate(ctx context.Context, prompts []string) (*golc.LLMResult, error) {
	subPromps := util.ChunkBy(prompts, o.opts.BatchSize)

	choices := []openai.CompletionChoice{}

	for _, prompt := range subPromps {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			res, err := o.client.CreateCompletion(ctx, openai.CompletionRequest{
				Prompt:      prompt,
				Model:       o.opts.Model,
				Temperature: o.opts.Temperatur,
				MaxTokens:   o.opts.MaxTokens,
				TopP:        o.opts.TopP,
			})
			if err != nil {
				return nil, err
			}

			choices = append(choices, res.Choices...)
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
		LLMOutput:   map[string]any{},
	}, nil
}
