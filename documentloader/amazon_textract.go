package documentloader

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/hupe1980/go-textractor"
	"github.com/hupe1980/golc/schema"
)

// AmazonTextractClient is an interface representing the methods required for interacting with Amazon Textract.
type AmazonTextractClient interface {
	// AnalyzeDocument performs document analysis using Amazon Textract.
	AnalyzeDocument(ctx context.Context, params *textract.AnalyzeDocumentInput, optFns ...func(*textract.Options)) (*textract.AnalyzeDocumentOutput, error)
}

// AmazonTextractOptions represents options for loading documents using Amazon Textract.
type AmazonTextractOptions struct {
	textractor.TextLinearizationOptions
	FeatureTypes []types.FeatureType
}

// DefaultLinerizationOptions returns the default linearization options for Amazon Textract.
func DefaultLinerizationOptions() textractor.TextLinearizationOptions {
	opts := textractor.DefaultLinerizationOptions
	opts.HideFigureLayout = true
	opts.TitlePrefix = "# "
	opts.SectionHeaderPrefix = "## "
	opts.ListElementPrefix = "* "

	return opts
}

// AmazonTextract represents a document loader for Amazon Textract.
type AmazonTextract struct {
	client   AmazonTextractClient
	r        io.Reader
	s3Object *types.S3Object
	output   *textractor.DocumentAPIOutput
	opts     AmazonTextractOptions
}

// NewAmazonTextractFromS3Object creates a new AmazonTextract instance from an S3 object.
func NewAmazonTextractFromS3Object(client AmazonTextractClient, s3Object *types.S3Object, optFns ...func(o *AmazonTextractOptions)) *AmazonTextract {
	opts := AmazonTextractOptions{
		TextLinearizationOptions: DefaultLinerizationOptions(),
		FeatureTypes:             []types.FeatureType{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &AmazonTextract{
		client:   client,
		s3Object: s3Object,
		opts:     opts,
	}
}

// NewAmazonTextractFromReader creates a new AmazonTextract instance from a reader.
func NewAmazonTextractFromReader(client AmazonTextractClient, r io.Reader, optFns ...func(o *AmazonTextractOptions)) *AmazonTextract {
	opts := AmazonTextractOptions{
		TextLinearizationOptions: DefaultLinerizationOptions(),
		FeatureTypes:             []types.FeatureType{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &AmazonTextract{
		client: client,
		r:      r,
		opts:   opts,
	}
}

// NewAmazonTextractFromOutput creates a new AmazonTextract instance from a Textract output.
func NewAmazonTextractFromOutput(output *textractor.DocumentAPIOutput, optFns ...func(o *AmazonTextractOptions)) *AmazonTextract {
	opts := AmazonTextractOptions{
		TextLinearizationOptions: DefaultLinerizationOptions(),
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &AmazonTextract{
		output: output,
		opts:   opts,
	}
}

// Load reads the content from the reader and returns it as a single document.
func (l *AmazonTextract) Load(ctx context.Context) ([]schema.Document, error) {
	var output *textractor.DocumentAPIOutput

	if l.output != nil {
		output = l.output
	} else if l.r != nil {
		o, err := l.analyzeBytes(ctx)
		if err != nil {
			return nil, err
		}

		output = &textractor.DocumentAPIOutput{
			DocumentMetadata: o.DocumentMetadata,
			Blocks:           o.Blocks,
		}
	} else if l.s3Object != nil {
		o, err := l.analyzeS3Object(ctx)
		if err != nil {
			return nil, err
		}

		output = &textractor.DocumentAPIOutput{
			DocumentMetadata: o.DocumentMetadata,
			Blocks:           o.Blocks,
		}
	} else {
		return nil, fmt.Errorf("unsupported api call")
	}

	doc, err := textractor.ParseDocumentAPIOutput(output)
	if err != nil {
		return nil, err
	}

	docs := make([]schema.Document, 0, len(doc.Pages()))

	for i, p := range doc.Pages() {
		docs = append(docs, schema.Document{
			PageContent: p.Text(func(tlo *textractor.TextLinearizationOptions) {
				*tlo = l.opts.TextLinearizationOptions
			}),
			Metadata: map[string]any{
				"page": i + 1,
			},
		})
	}

	return docs, nil
}

// analyzeBytes performs document analysis on the input bytes.
func (l *AmazonTextract) analyzeBytes(ctx context.Context) (*textract.AnalyzeDocumentOutput, error) {
	b, err := io.ReadAll(l.r)
	if err != nil {
		return nil, err
	}

	return l.client.AnalyzeDocument(ctx, &textract.AnalyzeDocumentInput{
		Document: &types.Document{
			Bytes: b,
		},
		FeatureTypes: l.opts.FeatureTypes,
	})
}

// analyzeS3Object performs document analysis on the S3 object.
func (l *AmazonTextract) analyzeS3Object(ctx context.Context) (*textract.AnalyzeDocumentOutput, error) {
	return l.client.AnalyzeDocument(ctx, &textract.AnalyzeDocumentInput{
		Document: &types.Document{
			S3Object: l.s3Object,
		},
		FeatureTypes: l.opts.FeatureTypes,
	})
}

// LoadAndSplit reads the content from the reader and splits it into multiple documents using the provided splitter.
func (l *AmazonTextract) LoadAndSplit(ctx context.Context, splitter schema.TextSplitter) ([]schema.Document, error) {
	docs, err := l.Load(ctx)
	if err != nil {
		return nil, err
	}

	return splitter.SplitDocuments(docs)
}
