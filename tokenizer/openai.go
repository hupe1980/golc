package tokenizer

import (
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

func (t *OpenAI) GetTokenIDs(text string) ([]uint, error) {
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

func (t *OpenAI) GetNumTokens(text string) (uint, error) {
	ids, err := t.GetTokenIDs(text)
	if err != nil {
		return 0, err
	}

	return uint(len(ids)), nil
}

func (t *OpenAI) GetNumTokensFromMessage(messages schema.ChatMessages) (uint, error) {
	text, err := messages.Format()
	if err != nil {
		return 0, err
	}

	return t.GetNumTokens(text)
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
