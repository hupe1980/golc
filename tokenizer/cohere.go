package tokenizer

import (
	"context"

	"github.com/cohere-ai/tokenizer"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Cohere satisfies the Tokenizer interface.
var _ schema.Tokenizer = (*Cohere)(nil)

type Cohere struct {
	encoder *tokenizer.Encoder
}

func NewCohere(modelName string) (*Cohere, error) {
	encoder, err := tokenizer.NewFromPrebuilt("coheretext-50k")
	if err != nil {
		return nil, err
	}

	return &Cohere{
		encoder: encoder,
	}, nil
}

// GetTokenIDs returns the token IDs corresponding to the provided text.
func (t *Cohere) GetTokenIDs(ctx context.Context, text string) ([]uint, error) {
	ids, _ := t.encoder.Encode(text)

	return int64ToUintSlice(ids), nil
}

// GetNumTokens returns the number of tokens in the provided text.
func (t *Cohere) GetNumTokens(ctx context.Context, text string) (uint, error) {
	ids, err := t.GetTokenIDs(ctx, text)
	if err != nil {
		return 0, err
	}

	return uint(len(ids)), nil
}

// GetNumTokensFromMessage returns the number of tokens in the provided chat messages.
func (t *Cohere) GetNumTokensFromMessage(ctx context.Context, messages schema.ChatMessages) (uint, error) {
	text, err := messages.Format()
	if err != nil {
		return 0, err
	}

	return t.GetNumTokens(ctx, text)
}

func int64ToUintSlice(numbers []int64) []uint {
	result := make([]uint, len(numbers))
	for i, num := range numbers {
		result[i] = uint(num)
	}

	return result
}
