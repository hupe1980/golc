package chatmodel

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
	"github.com/hupe1980/golc/util"
	"github.com/sashabaranov/go-openai"
)

// Compile time check to ensure OpenAI satisfies the ChatModel interface.
var _ schema.ChatModel = (*OpenAI)(nil)

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
		ModelName:        openai.GPT3Dot5Turbo,
		Temperatur:       1,
		TopP:             1,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Tokenizer == nil {
		opts.Tokenizer = tokenizer.NewOpenAI(opts.ModelName)
	}

	return newOpenAI(openai.NewClient(apiKey), opts)
}

func newOpenAI(client *openai.Client, opts OpenAIOptions) (*OpenAI, error) {
	return &OpenAI{
		Tokenizer: opts.Tokenizer,
		client:    client,
		opts:      opts,
	}, nil
}

func (cm *OpenAI) Generate(ctx context.Context, messages schema.ChatMessages, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	openAIMessages := []openai.ChatCompletionMessage{}

	for _, message := range messages {
		role, err := messageTypeToOpenAIRole(message.Type())
		if err != nil {
			return nil, err
		}

		if functionMessage, ok := message.(*schema.FunctionChatMessage); ok {
			openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
				Role:    role,
				Content: functionMessage.Content(),
				Name:    functionMessage.Name(),
			})
		} else {
			openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
				Role:    role,
				Content: message.Content(),
			})
		}
	}

	var functions []openai.FunctionDefinition
	if opts.Functions != nil {
		functions = util.Map(opts.Functions, func(fd schema.FunctionDefinition, i int) openai.FunctionDefinition {
			return openai.FunctionDefinition{
				Name:        fd.Name,
				Description: fd.Description,
				Parameters:  fd.Parameters,
			}
		})
	}

	res, err := cm.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:     cm.opts.ModelName,
		Messages:  openAIMessages,
		Functions: functions,
	})
	if err != nil {
		return nil, err
	}

	return &schema.ModelResult{
		Generations: []schema.Generation{{
			Text:    res.Choices[0].Message.Content,
			Message: openAIResponseToChatMessage(res.Choices[0].Message),
		}},
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
	case schema.ChatMessageTypeFunction:
		return "function", nil
	default:
		return "", fmt.Errorf("unknown message type: %s", mType)
	}
}

func openAIResponseToChatMessage(msg openai.ChatCompletionMessage) schema.ChatMessage {
	switch msg.Role {
	case "user":
		return schema.NewHumanChatMessage(msg.Content)
	case "assistant":
		if msg.FunctionCall != nil {
			return schema.NewAIChatMessage(msg.Content, func(o *schema.ChatMessageExtension) {
				o.FunctionCall = &schema.FunctionCall{
					Name:      msg.FunctionCall.Name,
					Arguments: msg.FunctionCall.Arguments,
				}
			})
		}

		return schema.NewAIChatMessage(msg.Content)
	case "system":
		return schema.NewSystemChatMessage(msg.Content)
	case "function":
		return schema.NewFunctionChatMessage(msg.Content, msg.Name)
	}

	return schema.NewGenericChatMessage(msg.Content, "unknown")
}

// Type returns the type of the model.
func (cm *OpenAI) Type() string {
	return "chatmodel.OpenAI"
}

// Verbose returns the verbosity setting of the model.
func (cm *OpenAI) Verbose() bool {
	return cm.opts.CallbackOptions.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (cm *OpenAI) Callbacks() []schema.Callback {
	return cm.opts.CallbackOptions.Callbacks
}

// InvocationParams returns the parameters used in the model invocation.
func (cm *OpenAI) InvocationParams() map[string]any {
	return nil
}
