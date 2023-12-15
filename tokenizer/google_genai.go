package tokenizer

import (
	"context"

	"cloud.google.com/go/ai/generativelanguage/apiv1/generativelanguagepb"
	"github.com/googleapis/gax-go/v2"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure GoogleGenAI satisfies the Tokenizer interface.
var _ schema.Tokenizer = (*GoogleGenAI)(nil)

// GoogleGenAIClient is an interface for the GoogleGenAI model client.
type GoogleGenAIClient interface {
	CountTokens(context.Context, *generativelanguagepb.CountTokensRequest, ...gax.CallOption) (*generativelanguagepb.CountTokensResponse, error)
}

type GoogleGenAI struct {
	client GoogleGenAIClient
	model  string
}

func NewGoogleGenAI(client GoogleGenAIClient, model string) *GoogleGenAI {
	return &GoogleGenAI{
		client: client,
		model:  model,
	}
}

// GetNumTokens returns the number of tokens in the provided text.
func (t *GoogleGenAI) GetNumTokens(ctx context.Context, text string) (uint, error) {
	res, err := t.client.CountTokens(ctx, &generativelanguagepb.CountTokensRequest{
		Model: t.model,
		Contents: []*generativelanguagepb.Content{{Parts: []*generativelanguagepb.Part{{
			Data: &generativelanguagepb.Part_Text{Text: text},
		}}}},
	})
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
