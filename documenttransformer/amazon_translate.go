package documenttransformer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/translate"
	"github.com/hupe1980/golc/schema"
	"golang.org/x/sync/errgroup"
)

// Compile time check to ensure AmazonTranslate satisfies the DocumentTransformer interface.
var _ schema.DocumentTransformer = (*AmazonTranslate)(nil)

// AmazonTranslateClient is an interface for the Amazon Translate client.
type AmazonTranslateClient interface {
	TranslateText(ctx context.Context, params *translate.TranslateTextInput, optFns ...func(*translate.Options)) (*translate.TranslateTextOutput, error)
}

// AmazonTranslateOptions contains options for configuring the AmazonTranslate transformer.
type AmazonTranslateOptions struct {
	// MaxConcurrency sets the maximum number of concurrent translation requests.
	MaxConcurrency int
	// SourceLanguageCode sets the source language code for translation. Default is "auto".
	SourceLanguageCode string
	// IncludeSourceText indicates whether to include the source text in document metadata.
	IncludeSourceText bool
}

// AmazonTranslate is a transformer that uses Amazon Translate to translate text in documents to a target language.
type AmazonTranslate struct {
	client             AmazonTranslateClient
	targetLanguageCode string
	opts               AmazonTranslateOptions
}

// NewAmazonTranslate creates a new instance of the AmazonTranslate transformer.
func NewAmazonTranslate(client AmazonTranslateClient, targetLanguageCode string, optFns ...func(o *AmazonTranslateOptions)) *AmazonTranslate {
	opts := AmazonTranslateOptions{
		MaxConcurrency:     5,
		SourceLanguageCode: "auto",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &AmazonTranslate{
		client:             client,
		targetLanguageCode: targetLanguageCode,
		opts:               opts,
	}
}

// Transform translates the text in a slice of documents to the target language using Amazon Translate.
func (t *AmazonTranslate) Transform(ctx context.Context, docs []schema.Document) ([]schema.Document, error) {
	errs, errctx := errgroup.WithContext(ctx)

	errs.SetLimit(t.opts.MaxConcurrency)

	translatedDocs := make([]schema.Document, len(docs))

	for i, d := range docs {
		i, d := i, d

		errs.Go(func() error {
			sourceText := d.PageContent

			if t.opts.SourceLanguageCode == t.targetLanguageCode {
				t.updateDocument(&d, sourceText, sourceText, t.opts.SourceLanguageCode, t.targetLanguageCode)
			}

			res, err := t.client.TranslateText(errctx, &translate.TranslateTextInput{
				SourceLanguageCode: aws.String(t.opts.SourceLanguageCode),
				TargetLanguageCode: aws.String(t.targetLanguageCode),
				Text:               aws.String(sourceText),
			})
			if err != nil {
				return err
			}

			t.updateDocument(&d, sourceText, aws.ToString(res.TranslatedText), aws.ToString(res.SourceLanguageCode), aws.ToString(res.TargetLanguageCode))

			translatedDocs[i] = d

			return nil
		})
	}

	if err := errs.Wait(); err != nil {
		return nil, err
	}

	return translatedDocs, nil
}

func (t *AmazonTranslate) updateDocument(doc *schema.Document, sourceText, translatedText, sourceLanguageCode, targetLanguageCode string) {
	doc.PageContent = translatedText

	if doc.Metadata == nil {
		doc.Metadata = make(map[string]any)
	}

	doc.Metadata["sourceLanguageCode"] = sourceLanguageCode
	doc.Metadata["targetLanguageCode"] = targetLanguageCode

	if t.opts.IncludeSourceText {
		doc.Metadata["sourceText"] = sourceText
	}
}
