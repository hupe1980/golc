package tokenizer

import "github.com/hupe1980/golc"

type Simple struct{}

func NewSimple() *Simple {
	return &Simple{}
}

func (t *Simple) GetTokenIDs(text string) ([]int, error) {
	return nil, nil
}

func (t *Simple) GetNumTokens(text string) (int, error) {
	return 0, nil
}

func (t *Simple) GetNumTokensFromMessage(messages []golc.ChatMessage) (int, error) {
	return 0, nil
}
