package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	bedrockruntimeTypes "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/integration/ai21"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
)

// Compile time check to ensure Bedrock satisfies the LLM interface.
var _ schema.LLM = (*Bedrock)(nil)

// providerStopSequenceKeyMap is a mapping between language model (LLM) providers
// and the corresponding key names used for stop sequences. Stop sequences are sets
// of words that, when encountered in the generated text, signal the language model
// to stop generating further content. Different LLM providers might use different
// key names to specify these stop sequences in the input parameters.
var providerStopSequenceKeyMap = map[string]string{
	"anthropic": "stop_sequences",
	"amazon":    "stopSequences",
	"ai21":      "stop_sequences",
	"cohere":    "stop_sequences",
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

		body["prompt"] = fmt.Sprintf("\n\nHuman:%s\n\nAssistant:", prompt)
	case "cohere":
		body = modelParams
		body["prompt"] = prompt
	case "meta":
		body = modelParams
		body["prompt"] = prompt
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

// cohereOutput is a struct representing the output structure for the "cohere" provider.
type cohereOutput struct {
	Generations []struct {
		Text string `json:"text"`
	} `json:"generations"`
}

// metaOutput is a struct representing the output structure for the "meta" provider.
type metaOutput struct {
	Generation string `json:"generation"`
}

// PrepareOutput prepares the output for the Bedrock model based on the specified provider.
func (bioa *BedrockInputOutputAdapter) PrepareOutput(response []byte) (string, error) {
	switch bioa.provider {
	case "ai21":
		output := &ai21Output{}
		if err := json.Unmarshal(response, output); err != nil {
			return "", err
		}

		return output.Completions[0].Data.Text, nil
	case "amazon":
		output := &amazonOutput{}
		if err := json.Unmarshal(response, output); err != nil {
			return "", err
		}

		return output.Results[0].OutputText, nil
	case "anthropic":
		output := &anthropicOutput{}
		if err := json.Unmarshal(response, output); err != nil {
			return "", err
		}

		return output.Completion, nil
	case "cohere":
		output := &cohereOutput{}
		if err := json.Unmarshal(response, output); err != nil {
			return "", err
		}

		return output.Generations[0].Text, nil
	case "meta":
		output := &metaOutput{}
		if err := json.Unmarshal(response, output); err != nil {
			return "", err
		}

		return output.Generation, nil
	}

	return "", fmt.Errorf("unsupported provider: %s", bioa.provider)
}

// BedrockRuntimeClient is an interface for the Bedrock model runtime client.
type BedrockRuntimeClient interface {
	InvokeModel(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error)
	InvokeModelWithResponseStream(ctx context.Context, params *bedrockruntime.InvokeModelWithResponseStreamInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelWithResponseStreamOutput, error)
}

type BedrockAI21Options struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// Model id to use.
	ModelID string `map:"model_id,omitempty"`

	// Temperature controls the randomness of text generation. Higher values make it more random.
	Temperature float64 `map:"temperature"`

	// TopP sets the nucleus sampling probability. Higher values result in more diverse text.
	TopP float64 `map:"topP"`

	// MaxTokens sets the maximum number of tokens in the generated text.
	MaxTokens int `map:"maxTokens"`

	// PresencePenalty specifies the penalty for repeating words in generated text.
	PresencePenalty ai21.Penalty `map:"presencePenalty"`

	// CountPenalty specifies the penalty for repeating tokens in generated text.
	CountPenalty ai21.Penalty `map:"countPenalty"`

	// FrequencyPenalty specifies the penalty for generating frequent words.
	FrequencyPenalty ai21.Penalty `map:"frequencyPenalty"`

	// Stream indicates whether to stream the results or not.
	Stream bool `map:"stream,omitempty"`
}

func NewBedrockAI21(client BedrockRuntimeClient, optFns ...func(o *BedrockAI21Options)) (*Bedrock, error) {
	opts := BedrockAI21Options{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelID:          "ai21.j2-ultra-v1", //https://docs.aws.amazon.com/bedrock/latest/userguide/model-ids-arns.html
		Temperature:      0.5,
		TopP:             0.5,
		MaxTokens:        200,
		PresencePenalty:  DefaultPenalty,
		CountPenalty:     DefaultPenalty,
		FrequencyPenalty: DefaultPenalty,
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

	return NewBedrock(client, func(o *BedrockOptions) {
		o.CallbackOptions = opts.CallbackOptions
		o.Tokenizer = opts.Tokenizer
		o.ModelID = opts.ModelID
		o.ModelParams = map[string]any{
			"temperature":      opts.Temperature,
			"topP":             opts.TopP,
			"maxTokens":        opts.MaxTokens,
			"presencePenalty":  opts.PresencePenalty,
			"countPenalty":     opts.CountPenalty,
			"frequencyPenalty": opts.FrequencyPenalty,
		}
		o.Stream = opts.Stream
	})
}

type BedrockAnthropicOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// Model id to use.
	ModelID string `map:"model_id,omitempty"`

	// MaxTokensToSmaple sets the maximum number of tokens in the generated text.
	MaxTokensToSample int `map:"max_tokens_to_sample"`

	// Temperature controls the randomness of text generation. Higher values make it more random.
	Temperature float32 `map:"temperature"`

	// TopP is the total probability mass of tokens to consider at each step.
	TopP float32 `map:"top_p,omitempty"`

	// TopK determines how the model selects tokens for output.
	TopK int `map:"top_k"`

	// Stream indicates whether to stream the results or not.
	Stream bool `map:"stream,omitempty"`
}

func NewBedrockAntrophic(client BedrockRuntimeClient, optFns ...func(o *BedrockAnthropicOptions)) (*Bedrock, error) {
	opts := BedrockAnthropicOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelID:           "anthropic.claude-v2", //https://docs.aws.amazon.com/bedrock/latest/userguide/model-ids-arns.html
		Temperature:       0.5,
		MaxTokensToSample: 256,
		TopP:              1,
		TopK:              250,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Tokenizer == nil {
		var tErr error

		opts.Tokenizer, tErr = tokenizer.NewClaude()
		if tErr != nil {
			return nil, tErr
		}
	}

	return NewBedrock(client, func(o *BedrockOptions) {
		o.CallbackOptions = opts.CallbackOptions
		o.Tokenizer = opts.Tokenizer
		o.ModelID = opts.ModelID
		o.ModelParams = map[string]any{
			"max_tokens_to_sample": opts.MaxTokensToSample,
			"temperature":          opts.Temperature,
			"top_p":                opts.TopP,
			"top_k":                opts.TopK,
		}
		o.Stream = opts.Stream
	})
}

type BedrockAmazonOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// Model id to use.
	ModelID string `map:"model_id,omitempty"`

	// Temperature controls the randomness of text generation. Higher values make it more random.
	Temperature float64 `json:"temperature"`

	// TopP is the total probability mass of tokens to consider at each step.
	TopP float64 `json:"topP"`

	// MaxTokenCount sets the maximum number of tokens in the generated text.
	MaxTokenCount int `json:"maxTokenCount"`

	// Stream indicates whether to stream the results or not.
	Stream bool `map:"stream,omitempty"`
}

func NewBedrockAmazon(client BedrockRuntimeClient, optFns ...func(o *BedrockAmazonOptions)) (*Bedrock, error) {
	opts := BedrockAmazonOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelID:       "amazon.titan-text-lite-v1", //https://docs.aws.amazon.com/bedrock/latest/userguide/model-ids-arns.html
		Temperature:   0,
		TopP:          1,
		MaxTokenCount: 512,
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

	return NewBedrock(client, func(o *BedrockOptions) {
		o.CallbackOptions = opts.CallbackOptions
		o.Tokenizer = opts.Tokenizer
		o.ModelID = opts.ModelID
		o.ModelParams = map[string]any{
			"temperature":   opts.Temperature,
			"topP":          opts.TopP,
			"maxTokenCount": opts.MaxTokenCount,
		}
		o.Stream = opts.Stream
	})
}

type ReturnLikelihood string

const (
	ReturnLikelihoodGeneration ReturnLikelihood = "GENERATION"
	ReturnLikelihoodAll        ReturnLikelihood = "ALL"
	ReturnLikelihoodNone       ReturnLikelihood = "NONE"
)

type BedrockCohereOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// Model id to use.
	ModelID string `map:"model_id,omitempty"`

	// Temperature controls the randomness of text generation. Higher values make it more random.
	Temperature float64 `json:"temperature,omitempty"`

	// P is the total probability mass of tokens to consider at each step.
	P float64 `json:"p,omitempty"`

	// K determines how the model selects tokens for output.
	K float64 `json:"k,omitempty"`

	// MaxTokens sets the maximum number of tokens in the generated text.
	MaxTokens int `json:"max_tokens,omitempty"`

	// ReturnLikelihoods specifies how and if the token likelihoods are returned with the response.
	ReturnLikelihoods ReturnLikelihood `json:"return_likelihoods,omitempty"`

	// Stream indicates whether to stream the results or not.
	Stream bool `map:"stream,omitempty"`
}

func NewBedrockCohere(client BedrockRuntimeClient, optFns ...func(o *BedrockCohereOptions)) (*Bedrock, error) {
	opts := BedrockCohereOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelID:           "cohere.command-text-v14", //https://docs.aws.amazon.com/bedrock/latest/userguide/model-ids-arns.html
		Temperature:       0.9,
		P:                 0.75,
		K:                 0,
		MaxTokens:         20,
		ReturnLikelihoods: ReturnLikelihoodNone,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Tokenizer == nil {
		var tErr error

		opts.Tokenizer, tErr = tokenizer.NewCohere(opts.ModelID)
		if tErr != nil {
			return nil, tErr
		}
	}

	return NewBedrock(client, func(o *BedrockOptions) {
		o.CallbackOptions = opts.CallbackOptions
		o.Tokenizer = opts.Tokenizer
		o.ModelID = opts.ModelID
		o.ModelParams = map[string]any{
			"temperature":        opts.Temperature,
			"p":                  opts.P,
			"k":                  opts.K,
			"max_tokens":         opts.MaxTokens,
			"return_likelihoods": opts.ReturnLikelihoods,
			"stream":             opts.Stream,
		}
		o.Stream = opts.Stream
	})
}

// BedrockMetaOptions contains options for configuring the Bedrock model with the "meta" provider.
type BedrockMetaOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// Model id to use.
	ModelID string `map:"model_id,omitempty"`

	// Temperature controls the randomness of text generation. Higher values make it more random.
	Temperature float32 `map:"temperature"`

	// TopP is the total probability mass of tokens to consider at each step.
	TopP float32 `map:"top_p,omitempty"`

	// MaxGenLen specify the maximum number of tokens to use in the generated response.
	MaxGenLen int `map:"max_gen_len"`

	// Stream indicates whether to stream the results or not.
	Stream bool `map:"stream,omitempty"`
}

// NewBedrockMeta creates a new instance of Bedrock for the "meta" provider.
func NewBedrockMeta(client BedrockRuntimeClient, optFns ...func(o *BedrockMetaOptions)) (*Bedrock, error) {
	opts := BedrockMetaOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelID:     "meta.llama2-70b-v1", //https://docs.aws.amazon.com/bedrock/latest/userguide/model-ids-arns.html
		Temperature: 0.5,
		TopP:        0.9,
		MaxGenLen:   512,
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

	return NewBedrock(client, func(o *BedrockOptions) {
		o.CallbackOptions = opts.CallbackOptions
		o.Tokenizer = opts.Tokenizer
		o.ModelID = opts.ModelID
		o.ModelParams = map[string]any{
			"temperature": opts.Temperature,
			"top_p":       opts.TopP,
			"max_gen_len": opts.MaxGenLen,
		}
		o.Stream = opts.Stream
	})
}

// BedrockOptions contains options for configuring the Bedrock LLM model.
type BedrockOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// Model id to use.
	ModelID string `map:"model_id,omitempty"`

	// Model params to use.
	ModelParams map[string]any `map:"model_params,omitempty"`

	// Stream indicates whether to stream the results or not.
	Stream bool `map:"stream,omitempty"`
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

	var completion string

	if l.opts.Stream {
		res, err := l.client.InvokeModelWithResponseStream(ctx, &bedrockruntime.InvokeModelWithResponseStreamInput{
			ModelId:     aws.String(l.opts.ModelID),
			Body:        body,
			Accept:      aws.String("application/json"),
			ContentType: aws.String("application/json"),
		})
		if err != nil {
			return nil, err
		}

		stream := res.GetStream()

		defer stream.Close()

		tokens := []string{}

		for event := range stream.Events() {
			switch v := event.(type) {
			case *bedrockruntimeTypes.ResponseStreamMemberChunk:
				token, err := bioa.PrepareOutput(v.Value.Bytes)
				if err != nil {
					return nil, err
				}

				if err := opts.CallbackManger.OnModelNewToken(ctx, &schema.ModelNewTokenManagerInput{
					Token: token,
				}); err != nil {
					return nil, err
				}

				tokens = append(tokens, token)
			}
		}

		completion = strings.Join(tokens, "")
	} else {
		res, err := l.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
			ModelId:     aws.String(l.opts.ModelID),
			Body:        body,
			Accept:      aws.String("application/json"),
			ContentType: aws.String("application/json"),
		})
		if err != nil {
			return nil, err
		}

		output, err := bioa.PrepareOutput(res.Body)
		if err != nil {
			return nil, err
		}

		completion = output
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
