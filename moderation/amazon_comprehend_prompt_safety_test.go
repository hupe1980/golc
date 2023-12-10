package moderation

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/comprehend/types"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestAmazonComprehendPromptSafety(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		inputText     string
		expectedError string
	}{
		{
			name:          "Moderation Passed",
			inputText:     "harmless content",
			expectedError: "",
		},
		{
			name:          "Moderation Failed",
			inputText:     "unsafe content",
			expectedError: "unsafe prompt detected",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctx := context.Background()

			score := float32(0.1)
			if strings.Contains(tc.inputText, "unsafe") {
				score = 0.9
			}

			fakeClient := &fakeAmazonComprehendPromptSafetyClient{
				response: &comprehend.ClassifyDocumentOutput{
					Classes: []types.DocumentClass{
						{
							Name:  aws.String("UNSAFE_PROMPT"),
							Score: aws.Float32(score),
						},
						{
							Name:  aws.String("SAFE_PROMPT"),
							Score: aws.Float32(1 - score),
						},
					},
				},
			}
			chain := NewAmazonComprehendPromptSafety(fakeClient)

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

type fakeAmazonComprehendPromptSafetyClient struct {
	response *comprehend.ClassifyDocumentOutput
	err      error
}

func (c *fakeAmazonComprehendPromptSafetyClient) ClassifyDocument(ctx context.Context, params *comprehend.ClassifyDocumentInput, optFns ...func(*comprehend.Options)) (*comprehend.ClassifyDocumentOutput, error) {
	return c.response, c.err
}

func (c *fakeAmazonComprehendPromptSafetyClient) Options() comprehend.Options {
	return comprehend.Options{
		Region: "us-east-1",
	}
}
