package tokenizer

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure GoogleGenAI satisfies the Tokenizer interface.
var _ schema.Tokenizer = (*GoogleGenAI)(nil)

type GoogleGenAIModel interface {
	CountTokens(ctx context.Context, parts ...genai.Part) (*genai.CountTokensResponse, error)
}

type GoogleGenAI struct {
	model GoogleGenAIModel
}

func NewGoogleGenAITokenizer(model GoogleGenAIModel) *GoogleGenAI {
	return &GoogleGenAI{
		model: model,
	}
}

// GetNumTokens returns the number of tokens in the provided text.
func (t *GoogleGenAI) GetNumTokens(ctx context.Context, text string) (uint, error) {
	res, err := t.model.CountTokens(ctx, genai.Text(text))
	if err != nil {
		return 0, err
	}

	return uint(res.TotalTokens), nil
}

// GetNumTokensFromMessage returns the number of tokens in the provided chat messages.
func (t *GoogleGenAI) GetNumTokensFromMessage(ctx context.Context, messages schema.ChatMessages) (uint, error) {
	text, err := messages.Format()
	if err != nil {
		return 0, err
	}

	return t.GetNumTokens(ctx, text)
}
