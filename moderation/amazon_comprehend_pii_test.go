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

func TestAmazonComprehendPII(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		inputText     string
		expectedError string
	}{
		{
			name:          "Moderation Passed",
			inputText:     "nopii",
			expectedError: "",
		},
		{
			name:          "Moderation Failed",
			inputText:     "pii",
			expectedError: "pii content found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctx := context.Background()

			score := float32(0.1)
			if tc.inputText == "pii" {
				score = 0.9
			}

			fakeClient := &fakeAmazonComprehendPIIClient{
				response: &comprehend.ContainsPiiEntitiesOutput{
					Labels: []types.EntityLabel{
						{Name: types.PiiEntityTypeName, Score: aws.Float32(score)},
					},
				},
			}
			chain := NewAmazonComprehendPII(fakeClient)

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

type fakeAmazonComprehendPIIClient struct {
	response *comprehend.ContainsPiiEntitiesOutput
	err      error
}

func (c *fakeAmazonComprehendPIIClient) ContainsPiiEntities(ctx context.Context, params *comprehend.ContainsPiiEntitiesInput, optFns ...func(*comprehend.Options)) (*comprehend.ContainsPiiEntitiesOutput, error) {
	return c.response, c.err
}
