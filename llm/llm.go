package llm

import (
	"github.com/hupe1980/golc"
)

type tokenizer struct{}

func (t *tokenizer) GetTokenIDs(text string) ([]int, error) {
	// TODO
	return nil, nil
}

func (t *tokenizer) GetNumTokens(text string) (int, error) {
	ids, err := t.GetTokenIDs(text)
	if err != nil {
		return 0, err
	}

	return len(ids), nil
}

func (t *tokenizer) GetNumTokensFromMessage(messages []golc.ChatMessage) (int, error) {
	text, err := golc.StringifyChatMessages(messages)
	if err != nil {
		return 0, err
	}

	return t.GetNumTokens(text)
}
