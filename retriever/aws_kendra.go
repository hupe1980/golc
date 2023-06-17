package retriever

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kendra"
	"github.com/aws/aws-sdk-go-v2/service/kendra/types"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure AWSKendra satisfies the Retriever interface.
var _ schema.Retriever = (*AWSKendra)(nil)

type KendraClient interface {
	Query(ctx context.Context, params *kendra.QueryInput, optFns ...func(*kendra.Options)) (*kendra.QueryOutput, error)
}

type AWSKendraOptions struct {
	// Number of documents to query for
	K int

	// Languagecode used for querying.
	LanguageCode string
}

type AWSKendra struct {
	client KendraClient
	index  string
	opts   AWSKendraOptions
}

func NewAWSKendra(client KendraClient, index string, optFns ...func(o *AWSKendraOptions)) *AWSKendra {
	opts := AWSKendraOptions{
		K:            3,
		LanguageCode: "en",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &AWSKendra{
		client: client,
		index:  index,
		opts:   opts,
	}
}

func (r *AWSKendra) GetRelevantDocuments(ctx context.Context, query string) ([]schema.Document, error) {
	return r.kendraQuery(ctx, query)
}

func (r *AWSKendra) kendraQuery(ctx context.Context, query string) ([]schema.Document, error) {
	out, err := r.client.Query(ctx, &kendra.QueryInput{
		IndexId:   aws.String(r.index),
		QueryText: aws.String(query),
		AttributeFilter: &types.AttributeFilter{
			AndAllFilters: []types.AttributeFilter{
				{
					EqualsTo: &types.DocumentAttribute{
						Key: aws.String("_language_code"),
						Value: &types.DocumentAttributeValue{
							StringValue: aws.String(r.opts.LanguageCode),
						},
					}},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	rCount := r.opts.K
	if len(out.ResultItems) < r.opts.K {
		rCount = len(out.ResultItems)
	}

	docs := []schema.Document{}

	for i, result := range out.ResultItems {
		if i > rCount {
			break
		}

		docTitle := aws.ToString(result.DocumentTitle.Text)
		docURI := aws.ToString(result.DocumentURI)
		docType := string(result.Type)

		text := aws.ToString(result.DocumentExcerpt.Text)
		if result.AdditionalAttributes != nil && aws.ToString(result.AdditionalAttributes[0].Key) == "AnswerText" {
			text = aws.ToString(result.AdditionalAttributes[0].Value.TextWithHighlightsValue.Text)
		}

		text = cleanResult(text)

		docs = append(docs, schema.Document{
			PageContent: fmt.Sprintf("Document Title: %s\nDocument Excerpt: %s\n", docTitle, text),
			Metadata: map[string]any{
				"source":  docURI,
				"title":   docTitle,
				"excerpt": text,
				"type":    docType,
			},
		})
	}

	return docs, nil
}

// cleanResult removes excess whitespace and ellipsis from the given string.
func cleanResult(resText string) string {
	cleanedText := regexp.MustCompile(`\s+`).ReplaceAllString(resText, " ")
	cleanedText = strings.ReplaceAll(cleanedText, "...", "")

	return cleanedText
}
