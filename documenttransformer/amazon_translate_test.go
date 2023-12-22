package documenttransformer

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/translate"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestAmazonTranslate(t *testing.T) {
	// Mocks
	client := &mockAmazonTranslateClient{
		TranslateTextFn: func(ctx context.Context, params *translate.TranslateTextInput, optFns ...func(*translate.Options)) (*translate.TranslateTextOutput, error) {
			return &translate.TranslateTextOutput{
				TranslatedText:     aws.String("TranslatedText"),
				SourceLanguageCode: aws.String("en"),
				TargetLanguageCode: aws.String("fr"),
			}, nil
		},
	}

	// Test data
	targetLanguage := "fr"
	options := AmazonTranslateOptions{
		MaxConcurrency:     2,
		SourceLanguageCode: "auto",
		IncludeSourceText:  true,
	}
	transformer := NewAmazonTranslate(client, targetLanguage, func(o *AmazonTranslateOptions) { *o = options })

	// Input documents
	docs := []schema.Document{
		{PageContent: "Hello", Metadata: map[string]any{"key": "value"}},
		{PageContent: "World", Metadata: map[string]any{"key": "value"}},
	}

	// Expected output documents
	expected := []schema.Document{
		{
			PageContent: "TranslatedText",
			Metadata: map[string]any{
				"sourceLanguageCode": "en",
				"targetLanguageCode": "fr",
				"sourceText":         "Hello",
				"key":                "value",
			},
		},
		{
			PageContent: "TranslatedText",
			Metadata: map[string]any{
				"sourceLanguageCode": "en",
				"targetLanguageCode": "fr",
				"sourceText":         "World",
				"key":                "value",
			},
		},
	}

	// Run the test
	transformed, err := transformer.Transform(context.Background(), docs)
	assert.NoError(t, err)
	assert.Equal(t, expected, transformed)
}

// mockAmazonTranslateClient is a mock implementation of the AmazonTranslateClient interface.
type mockAmazonTranslateClient struct {
	TranslateTextFn func(ctx context.Context, params *translate.TranslateTextInput, optFns ...func(*translate.Options)) (*translate.TranslateTextOutput, error)
}

func (m *mockAmazonTranslateClient) TranslateText(ctx context.Context, params *translate.TranslateTextInput, optFns ...func(*translate.Options)) (*translate.TranslateTextOutput, error) {
	return m.TranslateTextFn(ctx, params, optFns...)
}
