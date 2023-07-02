package llm

import (
	"context"
	"errors"
	"io"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
	"github.com/hupe1980/golc/util"
	"github.com/sashabaranov/go-openai"
)

// Compile time check to ensure OpenAI satisfies the LLM interface.
var _ schema.LLM = (*OpenAI)(nil)

type OpenAIOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`
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
	// Whether to stream the results or not.
	Stream bool
}

type OpenAI struct {
	schema.Tokenizer
	client *openai.Client
	opts   OpenAIOptions
}

func NewOpenAI(apiKey string, optFns ...func(o *OpenAIOptions)) (*OpenAI, error) {
	opts := OpenAIOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelName:        "text-davinci-002",
		Temperatur:       0.7,
		MaxTokens:        256,
		TopP:             1,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		N:                1,
		BatchSize:        20,
		Stream:           false,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Tokenizer == nil {
		opts.Tokenizer = tokenizer.NewOpenAI(opts.ModelName)
	}

	return &OpenAI{
		Tokenizer: opts.Tokenizer,
		client:    openai.NewClient(apiKey),
		opts:      opts,
	}, nil
}

func (l *OpenAI) Generate(ctx context.Context, prompts []string, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{}

	subPromps := util.ChunkBy(prompts, l.opts.BatchSize)

	choices := []openai.CompletionChoice{}
	tokenUsage := make(map[string]int)

	for _, prompt := range subPromps {
		completionRequest := openai.CompletionRequest{
			Prompt:      prompt,
			Model:       l.opts.ModelName,
			Temperature: l.opts.Temperatur,
			MaxTokens:   l.opts.MaxTokens,
			TopP:        l.opts.TopP,
			Stop:        opts.Stop,
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if l.opts.Stream {
				completionRequest.Stream = true

				stream, err := l.client.CreateCompletionStream(ctx, completionRequest)
				if err != nil {
					return nil, err
				}

				defer stream.Close()

			streamProcessing:
				for {
					select {
					case <-ctx.Done():
						return nil, ctx.Err()
					default:
						res, err := stream.Recv()
						if errors.Is(err, io.EOF) {
							break streamProcessing
						}

						if err != nil {
							return nil, err
						}

						if opts.CallbackManger != nil {
							if err := opts.CallbackManger.OnModelNewToken(res.Choices[0].Text); err != nil {
								return nil, err
							}
						}

						choices = append(choices, res.Choices...)
					}
				}
			} else {
				res, err := l.client.CreateCompletion(ctx, completionRequest)
				if err != nil {
					return nil, err
				}

				choices = append(choices, res.Choices...)
				tokenUsage["CompletionTokens"] += res.Usage.CompletionTokens
				tokenUsage["PromptTokens"] += res.Usage.PromptTokens
				tokenUsage["TotalTokens"] += res.Usage.TotalTokens
			}
		}
	}

	generations := util.Map(util.ChunkBy(choices, l.opts.N), func(promptChoices []openai.CompletionChoice, _ int) []schema.Generation {
		return util.Map(promptChoices, func(choice openai.CompletionChoice, _ int) schema.Generation {
			return schema.Generation{
				Text: choice.Text,
				Info: map[string]any{
					"FinishReason": choice.FinishReason,
					"Logprobs":     choice.LogProbs,
				},
			}
		})
	})

	return &schema.ModelResult{
		Generations: generations,
		LLMOutput: map[string]any{
			"ModelName":  l.opts.ModelName,
			"TokenUsage": tokenUsage,
		},
	}, nil
}

func (l *OpenAI) Type() string {
	return "llm.OpenAI"
}

func (l *OpenAI) Verbose() bool {
	return l.opts.CallbackOptions.Verbose
}

func (l *OpenAI) Callbacks() []schema.Callback {
	return l.opts.CallbackOptions.Callbacks
}

func (l *OpenAI) InvocationParams() map[string]any {
	return util.StructToMap(l.opts)
}
