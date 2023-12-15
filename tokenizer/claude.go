package tokenizer

import (
	"context"

	"github.com/hupe1980/go-tiktoken"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Claude satisfies the Tokenizer interface.
var _ schema.Tokenizer = (*Claude)(nil)

type Claude struct {
	encoding *tiktoken.Encoding
}

func NewClaude() (*Claude, error) {
	claude, err := tiktoken.NewClaude()
	if err != nil {
		return nil, err
	}

	encoding, err := tiktoken.NewEncoding(claude)
	if err != nil {
		return nil, err
	}

	return &Claude{
		encoding: encoding,
	}, nil
}

// GetTokenIDs returns the token IDs corresponding to the provided text.
func (t *Claude) GetTokenIDs(ctx context.Context, text string) ([]uint, error) {
	ids, _, err := t.encoding.Encode(text, nil, nil)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

// GetNumTokens returns the number of tokens in the provided text.
func (t *Claude) GetNumTokens(ctx context.Context, text string) (uint, error) {
	ids, err := t.GetTokenIDs(ctx, text)
	if err != nil {
		return 0, err
	}

	return uint(len(ids)), nil
}

// GetNumTokensFromMessage returns the number of tokens in the provided chat messages.
func (t *Claude) GetNumTokensFromMessage(ctx context.Context, messages schema.ChatMessages) (uint, error) {
	text, err := messages.Format()
	if err != nil {
		return 0, err
	}

	return t.GetNumTokens(ctx, text)
}
