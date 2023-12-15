package tokenizer

import (
	"context"

	"github.com/hupe1980/go-tiktoken"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure GPT2 satisfies the Tokenizer interface.
var _ schema.Tokenizer = (*GPT2)(nil)

type GPT2 struct {
	encoding *tiktoken.Encoding
}

func NewGPT2() (*GPT2, error) {
	gpt2, err := tiktoken.NewGPT2()
	if err != nil {
		return nil, err
	}

	encoding, err := tiktoken.NewEncoding(gpt2)
	if err != nil {
		return nil, err
	}

	return &GPT2{
		encoding: encoding,
	}, nil
}

// GetTokenIDs returns the token IDs corresponding to the provided text.
func (t *GPT2) GetTokenIDs(ctx context.Context, text string) ([]uint, error) {
	ids, _, err := t.encoding.Encode(text, nil, nil)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

// GetNumTokens returns the number of tokens in the provided text.
func (t *GPT2) GetNumTokens(ctx context.Context, text string) (uint, error) {
	ids, err := t.GetTokenIDs(ctx, text)
	if err != nil {
		return 0, err
	}

	return uint(len(ids)), nil
}

// GetNumTokensFromMessage returns the number of tokens in the provided chat messages.
func (t *GPT2) GetNumTokensFromMessage(ctx context.Context, messages schema.ChatMessages) (uint, error) {
	text, err := messages.Format()
	if err != nil {
		return 0, err
	}

	return t.GetNumTokens(ctx, text)
}
