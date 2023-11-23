package moderation

import (
	"context"
	"errors"
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/require"
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
			chain, err := NewOpenAIFromClient(fakeClient)
			require.NoError(t, err)

			// Test
			inputs := schema.ChainValues{
				"input": tc.inputText,
			}
			outputs, err := chain.Call(ctx, inputs)

			// Assertions
			if tc.expectedError == "" {
				require.NoError(t, err)
				require.NotNil(t, outputs)
			} else {
				require.Nil(t, outputs)
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}

	// Test case with a custom OpenAIModerateFunc
	t.Run("Custom OpenAIModerateFunc", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		fakeClient := &fakeOpenAIClient{
			response: openai.ModerationResponse{
				ID:      "12345",
				Model:   "text-moderation-latest",
				Results: []openai.Result{{Flagged: false}},
			},
		}

		chain, err := NewOpenAIFromClient(fakeClient, func(o *OpenAIOptions) {
			o.OpenAIModerateFunc = func(id, model string, result openai.Result) (schema.ChainValues, error) {
				if result.Flagged {
					return nil, errors.New("custom content policy violation")
				}

				return schema.ChainValues{
					"output": "Custom func result",
				}, nil
			}
		})
		require.NoError(t, err)

		// Test
		inputs := schema.ChainValues{
			"input": "Some input text",
		}
		outputs, err := chain.Call(ctx, inputs)

		// Assertions
		require.NoError(t, err)
		require.NotNil(t, outputs)
		require.Equal(t, "Custom func result", outputs["output"])
	})

	t.Run("Custom Content Policy Violation", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		fakeClient := &fakeOpenAIClient{
			response: openai.ModerationResponse{
				ID:      "12345",
				Model:   "text-moderation-latest",
				Results: []openai.Result{{Flagged: true}},
			},
		}

		chain, err := NewOpenAIFromClient(fakeClient, func(o *OpenAIOptions) {
			o.OpenAIModerateFunc = func(id, model string, result openai.Result) (schema.ChainValues, error) {
				if result.Flagged {
					return nil, errors.New("custom content policy violation")
				}

				return schema.ChainValues{
					"output": "Custom func result",
				}, nil
			}
		})
		require.NoError(t, err)

		// Test
		inputs := schema.ChainValues{
			"input": "Some flagged text",
		}
		_, err = chain.Call(ctx, inputs)

		// Assertions
		require.Error(t, err)
		require.EqualError(t, err, "custom content policy violation")
	})
}

type fakeOpenAIClient struct {
	response openai.ModerationResponse
	err      error
}

func (c *fakeOpenAIClient) Moderations(ctx context.Context, request openai.ModerationRequest) (openai.ModerationResponse, error) {
	return c.response, c.err
}
