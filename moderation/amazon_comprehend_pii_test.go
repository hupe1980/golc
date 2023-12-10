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

func TestAmazonComprehendPII(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		inputText     string
		redact        bool
		expectedError string
		expectedText  string
	}{
		{
			name:          "Moderation Passed",
			inputText:     "harmless content",
			redact:        false,
			expectedError: "",
			expectedText:  "harmless content",
		},
		{
			name:          "Moderation Failed",
			inputText:     "pii",
			redact:        false,
			expectedError: "pii content found",
			expectedText:  "",
		},
		{
			name:          "Redacted",
			inputText:     "hello pii",
			redact:        true,
			expectedError: "",
			expectedText:  "hello ***",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctx := context.Background()

			score := float32(0.1)
			if strings.Contains(tc.inputText, "pii") {
				score = 0.9
			}

			fakeClient := &fakeAmazonComprehendPIIClient{
				containsResponse: &comprehend.ContainsPiiEntitiesOutput{
					Labels: []types.EntityLabel{
						{Name: types.PiiEntityTypeName, Score: aws.Float32(score)},
					},
				},
				detectResponse: &comprehend.DetectPiiEntitiesOutput{
					Entities: []types.PiiEntity{
						{Type: types.PiiEntityTypeName, Score: aws.Float32(score), BeginOffset: aws.Int32(6), EndOffset: aws.Int32(9)},
					},
				},
			}
			chain := NewAmazonComprehendPII(fakeClient, func(o *AmazonComprehendPIIOptions) {
				o.Redact = tc.redact
			})

			// Test
			inputs := schema.ChainValues{
				"input": tc.inputText,
			}
			outputs, err := chain.Call(ctx, inputs)

			// Assertions
			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.NotNil(t, outputs)
				assert.Equal(t, tc.expectedText, outputs["output"])
			} else {
				assert.Nil(t, outputs)
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

type fakeAmazonComprehendPIIClient struct {
	containsResponse *comprehend.ContainsPiiEntitiesOutput
	detectResponse   *comprehend.DetectPiiEntitiesOutput
	err              error
}

func (c *fakeAmazonComprehendPIIClient) ContainsPiiEntities(ctx context.Context, params *comprehend.ContainsPiiEntitiesInput, optFns ...func(*comprehend.Options)) (*comprehend.ContainsPiiEntitiesOutput, error) {
	return c.containsResponse, c.err
}

func (c *fakeAmazonComprehendPIIClient) DetectPiiEntities(ctx context.Context, params *comprehend.DetectPiiEntitiesInput, optFns ...func(*comprehend.Options)) (*comprehend.DetectPiiEntitiesOutput, error) {
	return c.detectResponse, c.err
}
