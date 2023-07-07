package llm

import (
	"context"
	"fmt"

	huggingface "github.com/hupe1980/go-huggingface"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
)

// Compile time check to ensure HuggingFaceHub satisfies the LLM interface.
var _ schema.LLM = (*HuggingFaceHub)(nil)

type HuggingFaceHubClient interface {
	TextGeneration(ctx context.Context, req *huggingface.TextGenerationRequest) (huggingface.TextGenerationResponse, error)
	Text2TextGeneration(ctx context.Context, req *huggingface.Text2TextGenerationRequest) (huggingface.Text2TextGenerationResponse, error)
	Summarization(ctx context.Context, req *huggingface.SummarizationRequest) (huggingface.SummarizationResponse, error)
}

type HuggingFaceHubOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`
	Task                    string
}

type HuggingFaceHub struct {
	schema.Tokenizer
	client HuggingFaceHubClient
	opts   HuggingFaceHubOptions
}

func NewHuggingFaceHub(apiToken string, optFns ...func(o *HuggingFaceHubOptions)) (*HuggingFaceHub, error) {
	client := huggingface.NewInferenceClient(apiToken)
	return NewHuggingFaceHubFromClient(client, optFns...)
}

func NewHuggingFaceHubFromClient(client HuggingFaceHubClient, optFns ...func(o *HuggingFaceHubOptions)) (*HuggingFaceHub, error) {
	opts := HuggingFaceHubOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		Task: "text-generation",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Tokenizer == nil {
		var tErr error

		opts.Tokenizer, tErr = tokenizer.NewGPT2()
		if tErr != nil {
			return nil, tErr
		}
	}

	return &HuggingFaceHub{
		Tokenizer: opts.Tokenizer,
		client:    client,
		opts:      opts,
	}, nil
}

// Generate generates text based on the provided prompt and options.
func (l *HuggingFaceHub) Generate(ctx context.Context, prompt string, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	var (
		text string
		err  error
	)

	if l.opts.Task == "text-generation" {
		text, err = l.textGeneration(ctx, prompt)
	} else if l.opts.Task == "text2text-generation" {
		text, err = l.text2textGeneration(ctx, prompt)
	} else if l.opts.Task == "summarization" {
		text, err = l.summarization(ctx, prompt)
	} else {
		err = fmt.Errorf("unknown task: %s", l.opts.Task)
	}

	if err != nil {
		return nil, err
	}

	return &schema.ModelResult{
		Generations: []schema.Generation{{Text: text}},
		LLMOutput:   map[string]any{},
	}, nil
}

func (l *HuggingFaceHub) textGeneration(ctx context.Context, input string) (string, error) {
	res, err := l.client.TextGeneration(ctx, &huggingface.TextGenerationRequest{
		Inputs: input,
	})
	if err != nil {
		return "", err
	}

	// Text generation return includes the starter text.
	return res[0].GeneratedText[len(input):], nil
}

func (l *HuggingFaceHub) text2textGeneration(ctx context.Context, input string) (string, error) {
	res, err := l.client.Text2TextGeneration(ctx, &huggingface.Text2TextGenerationRequest{
		Inputs: input,
	})
	if err != nil {
		return "", err
	}

	return res[0].GeneratedText, nil
}

func (l *HuggingFaceHub) summarization(ctx context.Context, input string) (string, error) {
	res, err := l.client.Summarization(ctx, &huggingface.SummarizationRequest{
		Inputs: []string{input},
	})
	if err != nil {
		return "", err
	}

	return res[0].SummaryText, nil
}

// Type returns the type of the model.
func (l *HuggingFaceHub) Type() string {
	return "llm.HuggingFaceHub"
}

// Verbose returns the verbosity setting of the model.
func (l *HuggingFaceHub) Verbose() bool {
	return l.opts.CallbackOptions.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (l *HuggingFaceHub) Callbacks() []schema.Callback {
	return l.opts.CallbackOptions.Callbacks
}

// InvocationParams returns the parameters used in the model invocation.
func (l *HuggingFaceHub) InvocationParams() map[string]any {
	return nil
}
