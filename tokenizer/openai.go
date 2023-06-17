package tokenizer

import (
	"github.com/hupe1980/golc"
	"github.com/pkoukk/tiktoken-go"
)

type OpenAI struct {
	modelName string
}

func NewOpenAI(modelName string) *OpenAI {
	return &OpenAI{
		modelName: modelName,
	}
}

func (o *OpenAI) GetTokenIDs(text string) ([]int, error) {
	_, e, err := o.getEncodingForModel()
	if err != nil {
		return nil, err
	}

	return e.Encode(text, nil, nil), nil
}

func (o *OpenAI) GetNumTokens(text string) (int, error) {
	ids, err := o.GetTokenIDs(text)
	if err != nil {
		return 0, err
	}

	return len(ids), nil
}

func (o *OpenAI) GetNumTokensFromMessage(messages []golc.ChatMessage) (int, error) {
	text, err := golc.StringifyChatMessages(messages)
	if err != nil {
		return 0, err
	}

	return o.GetNumTokens(text)
}

func (o *OpenAI) getEncodingForModel() (string, *tiktoken.Tiktoken, error) {
	model := o.modelName
	if model == "gpt-3.5-turbo" {
		model = "gpt-3.5-turbo-0301"
	} else if model == "gpt-4" {
		model = "gpt-4-0314"
	}

	e, err := tiktoken.EncodingForModel(model)
	if err != nil {
		model = "cl100k_base" //fallback

		e, err = tiktoken.EncodingForModel(model)

		return model, e, err
	}

	return model, e, nil
}
