package retriever

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kendra"
	"github.com/aws/aws-sdk-go-v2/service/kendra/types"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure AmazonKendra satisfies the Retriever interface.
var _ schema.Retriever = (*AmazonKendra)(nil)

// AmazonKendraClient represents a client for interacting with Amazon Kendra.
type AmazonKendraClient interface {
	// Retrieve retrieves documents from Amazon Kendra based on the provided input parameters.
	// It returns the retrieval output or an error if the retrieval operation fails.
	Retrieve(ctx context.Context, params *kendra.RetrieveInput, optFns ...func(*kendra.Options)) (*kendra.RetrieveOutput, error)

	// Query performs a query operation on Amazon Kendra using the specified input parameters.
	// It returns the query output or an error if the query operation fails.
	Query(ctx context.Context, params *kendra.QueryInput, optFns ...func(*kendra.Options)) (*kendra.QueryOutput, error)
}

type AmazonKendraOptions struct {
	*schema.CallbackOptions
	// Number of documents to query for
	TopK int32

	// Provides filtering the results based on document attributes or metadata
	// fields.
	AttributeFilter *types.AttributeFilter
}

type AmazonKendra struct {
	client AmazonKendraClient
	index  string
	opts   AmazonKendraOptions
}

func NewAmazonKendra(client AmazonKendraClient, index string, optFns ...func(o *AmazonKendraOptions)) *AmazonKendra {
	opts := AmazonKendraOptions{
		TopK: 3,
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &AmazonKendra{
		client: client,
		index:  index,
		opts:   opts,
	}
}

func (r *AmazonKendra) GetRelevantDocuments(ctx context.Context, query string) ([]schema.Document, error) {
	return r.kendraQuery(ctx, query)
}

// Verbose returns the verbosity setting of the retriever.
func (r *AmazonKendra) Verbose() bool {
	return r.opts.CallbackOptions.Verbose
}

// Callbacks returns the registered callbacks of the retriever.
func (r *AmazonKendra) Callbacks() []schema.Callback {
	return r.opts.CallbackOptions.Callbacks
}

func (r *AmazonKendra) kendraQuery(ctx context.Context, query string) ([]schema.Document, error) {
	query = strings.TrimSpace(query)

	retrieveOutput, err := r.client.Retrieve(ctx, &kendra.RetrieveInput{
		IndexId:         aws.String(r.index),
		QueryText:       aws.String(query),
		PageSize:        aws.Int32(r.opts.TopK),
		AttributeFilter: r.opts.AttributeFilter,
	})
	if err != nil {
		return nil, err
	}

	docs := []schema.Document{}
	for _, item := range retrieveOutput.ResultItems {
		docs = append(docs, r.parseRetrievalResultItem(item))
	}

	if len(retrieveOutput.ResultItems) == 0 {
		queryOutput, err := r.client.Query(ctx, &kendra.QueryInput{
			IndexId:         aws.String(r.index),
			QueryText:       aws.String(query),
			PageSize:        aws.Int32(r.opts.TopK),
			AttributeFilter: r.opts.AttributeFilter,
		})
		if err != nil {
			return nil, err
		}

		for _, item := range queryOutput.ResultItems {
			docs = append(docs, r.parseQueryResultItem(item))
		}
	}

	return docs, nil
}

func (r *AmazonKendra) parseRetrievalResultItem(item types.RetrieveResultItem) schema.Document {
	title := aws.ToString(item.DocumentTitle)
	content := r.cleanText(aws.ToString(item.Content))
	source := aws.ToString(item.DocumentURI)

	return schema.Document{
		PageContent: fmt.Sprintf("Document Title: %s\nDocument Excerpt: %s\n", title, content),
		Metadata: map[string]any{
			"source":  source,
			"title":   title,
			"excerpt": content,
		},
	}
}

func (r *AmazonKendra) parseQueryResultItem(item types.QueryResultItem) schema.Document {
	var text string
	if item.DocumentExcerpt != nil {
		text = aws.ToString(item.DocumentExcerpt.Text)
	}

	if len(item.AdditionalAttributes) > 0 && aws.ToString(item.AdditionalAttributes[0].Key) == "AnswerText" {
		text = aws.ToString(item.AdditionalAttributes[0].Value.TextWithHighlightsValue.Text)
	}

	text = r.cleanText(text)

	var title string
	if item.DocumentTitle != nil {
		title = aws.ToString(item.DocumentTitle.Text)
	}

	source := aws.ToString(item.DocumentURI)

	dtype := string(item.Type)

	return schema.Document{
		PageContent: fmt.Sprintf("Document Title: %s\nDocument Excerpt: %s\n", title, text),
		Metadata: map[string]any{
			"source":  source,
			"title":   title,
			"excerpt": text,
			"type":    dtype,
		},
	}
}

// cleanText removes excess whitespace and ellipsis from the given string.
func (r *AmazonKendra) cleanText(resText string) string {
	if resText == "" {
		return ""
	}

	cleanedText := regexp.MustCompile(`\s+`).ReplaceAllString(resText, " ")
	cleanedText = strings.ReplaceAll(cleanedText, "...", "")

	return cleanedText
}

func AmazonKendraLanguageCodeAttributeFilter(languageCode string) *types.AttributeFilter {
	return &types.AttributeFilter{
		AndAllFilters: []types.AttributeFilter{
			{
				EqualsTo: &types.DocumentAttribute{
					Key: aws.String("_language_code"),
					Value: &types.DocumentAttributeValue{
						StringValue: aws.String(languageCode),
					},
				}},
		},
	}
}
