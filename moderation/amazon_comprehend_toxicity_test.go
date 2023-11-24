package moderation

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/comprehend/types"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestAmazonComprehendToxicity(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		inputText     string
		expectedError string
	}{
		{
			name:          "Moderation Passed",
			inputText:     "notoxic",
			expectedError: "",
		},
		{
			name:          "Moderation Failed",
			inputText:     "toxic",
			expectedError: "toxic content found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctx := context.Background()

			score := float32(0.1)
			if tc.inputText == "toxic" {
				score = 0.9
			}

			fakeClient := &fakeAmazonComprehendToxicityClient{
				response: &comprehend.DetectToxicContentOutput{
					ResultList: []types.ToxicLabels{
						{
							Toxicity: aws.Float32(score),
							Labels: []types.ToxicContent{
								{
									Name:  types.ToxicContentTypeHateSpeech,
									Score: aws.Float32(score),
								},
							},
						},
					},
				},
			}
			chain := NewAmazonComprehendToxicity(fakeClient)

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

type fakeAmazonComprehendToxicityClient struct {
	response *comprehend.DetectToxicContentOutput
	err      error
}

func (c *fakeAmazonComprehendToxicityClient) DetectToxicContent(ctx context.Context, params *comprehend.DetectToxicContentInput, optFns ...func(*comprehend.Options)) (*comprehend.DetectToxicContentOutput, error) {
	return c.response, c.err
}
