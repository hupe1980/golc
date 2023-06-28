package llm

import (
	"context"
	"fmt"

	huggingface "github.com/hupe1980/go-huggingface"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
)

// Compile time check to ensure HuggingFaceHub satisfies the LLM interface.
var _ schema.LLM = (*HuggingFaceHub)(nil)

type HuggingFaceHubOptions struct {
	*schema.CallbackOptions
	Tokenizer schema.Tokenizer
	huggingface.HTTPClient
	// Model name to use.
	Model string
	Task  string
}

type HuggingFaceHub struct {
	schema.Tokenizer
	client *huggingface.InferenceClient
	opts   HuggingFaceHubOptions
}

func NewHuggingFaceHub(apiToken string, optFns ...func(o *HuggingFaceHubOptions)) (*HuggingFaceHub, error) {
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
		client: huggingface.NewInferenceClient(apiToken, func(o *huggingface.InferenceClientOptions) {
			o.Model = opts.Model
			o.HTTPClient = opts.HTTPClient
		}),
		opts: opts,
	}, nil
}

func (l *HuggingFaceHub) Generate(ctx context.Context, prompts []string, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	var (
		text string
		err  error
	)

	if l.opts.Task == "text-generation" {
		text, err = l.textGeneration(ctx, prompts[0])
	} else if l.opts.Task == "text2text-generation" {
		text, err = l.text2textGeneration(ctx, prompts[0])
	} else if l.opts.Task == "summarization" {
		text, err = l.summarization(ctx, prompts[0])
	} else {
		err = fmt.Errorf("unknown task: %s", l.opts.Task)
	}

	if err != nil {
		return nil, err
	}

	return &schema.ModelResult{
		Generations: [][]schema.Generation{{schema.Generation{Text: text}}},
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

func (l *HuggingFaceHub) Type() string {
	return "HuggingFaceHub"
}

func (l *HuggingFaceHub) Verbose() bool {
	return l.opts.CallbackOptions.Verbose
}

func (l *HuggingFaceHub) Callbacks() []schema.Callback {
	return l.opts.CallbackOptions.Callbacks
}
