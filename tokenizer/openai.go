package tokenizer

import (
	"context"
	"fmt"
	"strings"

	"github.com/hupe1980/go-tiktoken"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure OpenAI satisfies the Tokenizer interface.
var _ schema.Tokenizer = (*OpenAI)(nil)

type OpenAI struct {
	modelName string
}

func NewOpenAI(modelName string) *OpenAI {
	return &OpenAI{
		modelName: modelName,
	}
}

// GetTokenIDs returns the token IDs corresponding to the provided text.
func (t *OpenAI) GetTokenIDs(ctx context.Context, text string) ([]uint, error) {
	_, e, err := t.getEncodingForModel()
	if err != nil {
		return nil, err
	}

	ids, _, err := e.Encode(text, nil, nil)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

// GetNumTokens returns the number of tokens in the provided text.
func (t *OpenAI) GetNumTokens(ctx context.Context, text string) (uint, error) {
	ids, err := t.GetTokenIDs(ctx, text)
	if err != nil {
		return 0, err
	}

	return uint(len(ids)), nil
}

// GetNumTokensFromMessage returns the number of tokens in the provided chat messages.
func (t *OpenAI) GetNumTokensFromMessage(ctx context.Context, messages schema.ChatMessages) (uint, error) {
	var tokensPerMessage, tokensPerName int

	// Official documentation: https://github.com/openai/openai-cookbook/blob/main/examples/How_to_format_inputs_to_ChatGPT_models.ipynb"""
	if strings.HasPrefix(t.modelName, "gpt-3.5-turbo-0301") {
		// every message follows <im_start>{role/name}\n{content}<im_end>\n
		tokensPerMessage = 4
		// if there's a name, the role is omitted
		tokensPerName = -1
	} else if strings.HasPrefix(t.modelName, "gpt-3.5-turbo") || strings.HasPrefix(t.modelName, "gpt-4") {
		tokensPerMessage = 3
		tokensPerName = 1
	} else {
		return 0, fmt.Errorf("unsupported model: %s", t.modelName)
	}

	var numTokens int
	for _, m := range messages {
		numTokens += tokensPerMessage

		if m.Type() == schema.ChatMessageTypeFunction {
			fm, _ := m.(schema.FunctionChatMessage)
			if fm.Name() != "" {
				numTokens += tokensPerName
			}
		}

		nt, err := t.GetNumTokens(ctx, m.Content())
		if err != nil {
			return 0, err
		}

		numTokens += int(nt)
	}

	return uint(numTokens), nil
}

func (t *OpenAI) getEncodingForModel() (string, *tiktoken.Encoding, error) {
	model := t.modelName
	if model == "gpt-3.5-turbo" {
		model = "gpt-3.5-turbo-0301"
	} else if model == "gpt-4" {
		model = "gpt-4-0314"
	}

	e, err := tiktoken.NewEncodingForModel(model)
	if err != nil {
		model = "cl100k_base" //fallback

		e, err = tiktoken.NewEncodingForModel(model)

		return model, e, err
	}

	return model, e, nil
}
