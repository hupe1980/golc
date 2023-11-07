package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
	"github.com/hupe1980/golc/util"
)

// Compile time check to ensure Bedrock satisfies the LLM interface.
var _ schema.LLM = (*Bedrock)(nil)

const (
	humanPrompt     = "\n\nHuman:"
	assistantPrompt = "\n\nAssistant:"
)

func humanAssistantFormat(inputText string) string {
	inputText = fmt.Sprintf("%s %s", humanPrompt, inputText)
	if strings.Count(inputText, "Assistant:") == 0 {
		inputText = fmt.Sprintf("%s%s", inputText, assistantPrompt)
	}

	return inputText
}

// providerStopSequenceKeyMap is a mapping between language model (LLM) providers
// and the corresponding key names used for stop sequences. Stop sequences are sets
// of words that, when encountered in the generated text, signal the language model
// to stop generating further content. Different LLM providers might use different
// key names to specify these stop sequences in the input parameters.
var providerStopSequenceKeyMap = map[string]string{
	"anthropic": "stop_sequences",
	"amazon":    "stopSequences",
	"ai21":      "stop_sequences",
}

// BedrockInputOutputAdapter is a helper struct for preparing input and handling output for Bedrock model.
type BedrockInputOutputAdapter struct {
	provider string
}

// NewBedrockInputOutputAdpter creates a new instance of BedrockInputOutputAdpter.
func NewBedrockInputOutputAdapter(provider string) *BedrockInputOutputAdapter {
	return &BedrockInputOutputAdapter{
		provider: provider,
	}
}

// PrepareInput prepares the input for the Bedrock model based on the specified provider.
func (bioa *BedrockInputOutputAdapter) PrepareInput(prompt string, modelParams map[string]any) ([]byte, error) {
	var body map[string]any

	switch bioa.provider {
	case "ai21":
		body = modelParams
		body["prompt"] = prompt
	case "amazon":
		body = make(map[string]any)
		body["inputText"] = prompt
		body["textGenerationConfig"] = modelParams
	case "anthropic":
		body = modelParams

		if _, ok := body["max_tokens_to_sample"]; !ok {
			body["max_tokens_to_sample"] = 256
		}

		body["prompt"] = humanAssistantFormat(prompt)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", bioa.provider)
	}

	return json.Marshal(body)
}

// ai21Output represents the structure of the output from the AI21 language model.
// It is used for unmarshaling JSON responses from the language model's API.
type ai21Output struct {
	Completions []struct {
		Data struct {
			Text string `json:"text"`
		} `json:"data"`
	} `json:"completions"`
}

// amazonOutput represents the structure of the output from the Amazon language model.
// It is used for unmarshaling JSON responses from the language model's API.
type amazonOutput struct {
	Results []struct {
		OutputText string `json:"outputText"`
	} `json:"results"`
}

// anthropicOutput is a struct representing the output structure for the "anthropic" provider.
type anthropicOutput struct {
	Completion string `json:"completion"`
}

// PrepareOutput prepares the output for the Bedrock model based on the specified provider.
func (bioa *BedrockInputOutputAdapter) PrepareOutput(response *bedrockruntime.InvokeModelOutput) (string, error) {
	switch bioa.provider {
	case "ai21":
		output := &ai21Output{}
		if err := json.Unmarshal(response.Body, output); err != nil {
			return "", err
		}

		return output.Completions[0].Data.Text, nil
	case "amazon":
		output := &amazonOutput{}
		if err := json.Unmarshal(response.Body, output); err != nil {
			return "", err
		}

		return output.Results[0].OutputText, nil
	case "anthropic":
		output := &anthropicOutput{}
		if err := json.Unmarshal(response.Body, output); err != nil {
			return "", err
		}

		return output.Completion, nil
	}

	return "", fmt.Errorf("unsupported provider: %s", bioa.provider)
}

// BedrockRuntimeClient is an interface for the Bedrock model runtime client.
type BedrockRuntimeClient interface {
	InvokeModel(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error)
}

// BedrockOptions contains options for configuring the Bedrock LLM model.
type BedrockOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// Model id to use.
	ModelID string `map:"model_id,omitempty"`

	// Model params to use.
	ModelParams map[string]any `map:"model_params,omitempty"`
}

// Bedrock is a Bedrock LLM model that generates text based on a provided response function.
type Bedrock struct {
	schema.Tokenizer
	client BedrockRuntimeClient
	opts   BedrockOptions
}

// NewBedrock creates a new instance of the Bedrock LLM model with the provided response function and options.
func NewBedrock(client BedrockRuntimeClient, optFns ...func(o *BedrockOptions)) (*Bedrock, error) {
	opts := BedrockOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelParams: make(map[string]any),
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

	return &Bedrock{
		Tokenizer: opts.Tokenizer,
		client:    client,
		opts:      opts,
	}, nil
}

// Generate generates text based on the provided prompt and options.
func (l *Bedrock) Generate(ctx context.Context, prompt string, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	provider := l.getProvider()

	params := util.CopyMap(l.opts.ModelParams)

	if len(opts.Stop) > 0 {
		key, ok := providerStopSequenceKeyMap[provider]
		if !ok {
			return nil, fmt.Errorf("stop sequence key name for provider %s is not supported", provider)
		}

		params[key] = opts.Stop
	}

	bioa := NewBedrockInputOutputAdapter(provider)

	body, err := bioa.PrepareInput(prompt, params)
	if err != nil {
		return nil, err
	}

	res, err := l.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(l.opts.ModelID),
		Body:        body,
		Accept:      aws.String("application/json"),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return nil, err
	}

	completion, err := bioa.PrepareOutput(res)
	if err != nil {
		return nil, err
	}

	return &schema.ModelResult{
		Generations: []schema.Generation{{Text: completion}},
		LLMOutput:   map[string]any{},
	}, nil
}

// Type returns the type of the model.
func (l *Bedrock) Type() string {
	return "llm.Bedrock"
}

// Verbose returns the verbosity setting of the model.
func (l *Bedrock) Verbose() bool {
	return l.opts.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (l *Bedrock) Callbacks() []schema.Callback {
	return l.opts.Callbacks
}

// InvocationParams returns the parameters used in the model invocation.
func (l *Bedrock) InvocationParams() map[string]any {
	return util.StructToMap(l.opts)
}

// getProvider returns the provider of the model based on the model ID.
func (l *Bedrock) getProvider() string {
	return strings.Split(l.opts.ModelID, ".")[0]
}
