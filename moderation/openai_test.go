package moderation

import (
	"context"
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

func TestOpenAI(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		inputText     string
		flagged       bool
		expectedError string
	}{
		{
			name:          "Moderation Passed",
			inputText:     "Some text to moderate",
			flagged:       false,
			expectedError: "",
		},
		{
			name:          "Moderation Failed",
			inputText:     "Some flagged text",
			flagged:       true,
			expectedError: "content policy violation",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctx := context.Background()
			fakeClient := &fakeOpenAIClient{
				response: openai.ModerationResponse{
					ID:      "12345",
					Model:   "text-moderation-latest",
					Results: []openai.Result{{Flagged: tc.flagged}},
				},
			}
			chain := NewOpenAIFromClient(fakeClient)

			// Test
			inputs := schema.ChainValues{
				"input": tc.inputText,
			}
			outputs, err := chain.Call(ctx, inputs)

			// Assertions
			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.NotNil(t, outputs)
				assert.Equal(t, tc.inputText, outputs["output"])
			} else {
				assert.Nil(t, outputs)
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

type fakeOpenAIClient struct {
	response openai.ModerationResponse
	err      error
}

func (c *fakeOpenAIClient) Moderations(ctx context.Context, request openai.ModerationRequest) (openai.ModerationResponse, error) {
	return c.response, c.err
}
