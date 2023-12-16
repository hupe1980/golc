package documenttransformer

import (
	"context"
	"testing"

	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestMetaDataTagger(t *testing.T) {
	fake := chatmodel.NewFake(func(ctx context.Context, messages schema.ChatMessages) (*schema.ModelResult, error) {
		return &schema.ModelResult{
			Generations: []schema.Generation{{
				Message: schema.NewAIChatMessage("", func(o *schema.ChatMessageExtension) {
					o.FunctionCall = &schema.FunctionCall{
						Name:      "InformationExtraction",
						Arguments: "{\"foo\":\"bar\"}",
					}
				}),
			}},
		}, nil
	})

	// Mock data
	mockDocs := []schema.Document{
		{PageContent: "Content1"},
		{PageContent: "Content2"},
	}

	type tags struct {
		Foo string `json:"foo"`
	}

	// Create a new MetaDataTagger with a MockChatModel
	tagger, err := NewMetaDataTagger(fake, &tags{})
	assert.NoError(t, err)

	// Run the transformation
	result, err := tagger.Transform(context.Background(), mockDocs)
	assert.NoError(t, err)

	// Assertions
	assert.Len(t, result, len(mockDocs))

	for i, doc := range result {
		assert.Equal(t, mockDocs[i].PageContent, doc.PageContent)
		assert.Contains(t, doc.Metadata, "foo")
	}
}
