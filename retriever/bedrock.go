package retriever

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure BedrockKnowledgeBases satisfies the Retriever interface.
var _ schema.Retriever = (*BedrockKnowledgeBases)(nil)

// BedrockAgentRuntimeClient is an interface representing the Bedrock Agent Runtime client.
type BedrockAgentRuntimeClient interface {
	Retrieve(context.Context, *bedrockagentruntime.RetrieveInput, ...func(*bedrockagentruntime.Options)) (*bedrockagentruntime.RetrieveOutput, error)
}

// BedrockKnowledgeBasesOptions represents the options for configuring BedrockKnowledgeBases.
type BedrockKnowledgeBasesOptions struct {
	*schema.CallbackOptions

	// RetrievalConfiguration provides search parameters for retrieving from knowledge base.
	RetrievalConfiguration types.KnowledgeBaseRetrievalConfiguration
}

// BedrockKnowledgeBases is a retriever implementation for retrieving documents from a knowledge base using the Bedrock Agent Runtime client.
type BedrockKnowledgeBases struct {
	client          BedrockAgentRuntimeClient
	knowledgeBaseID string
	opts            BedrockKnowledgeBasesOptions
}

// NewBedrockKnowledgeBases creates a new BedrockKnowledgeBases retriever with the specified client, knowledge base ID, and options.
func NewBedrockKnowledgeBases(client BedrockAgentRuntimeClient, knowledgeBaseID string, optFns ...func(o *BedrockKnowledgeBasesOptions)) *BedrockKnowledgeBases {
	opts := BedrockKnowledgeBasesOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		RetrievalConfiguration: types.KnowledgeBaseRetrievalConfiguration{
			VectorSearchConfiguration: &types.KnowledgeBaseVectorSearchConfiguration{
				NumberOfResults: aws.Int32(3),
			},
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &BedrockKnowledgeBases{
		client:          client,
		knowledgeBaseID: knowledgeBaseID,
		opts:            opts,
	}
}

// GetRelevantDocuments retrieves relevant documents from the knowledge base based on the given query.
func (r *BedrockKnowledgeBases) GetRelevantDocuments(ctx context.Context, query string) ([]schema.Document, error) {
	query = strings.TrimSpace(query)

	docs := []schema.Document{}

	p := bedrockagentruntime.NewRetrievePaginator(r.client, &bedrockagentruntime.RetrieveInput{
		KnowledgeBaseId: aws.String(r.knowledgeBaseID),
		RetrievalQuery: &types.KnowledgeBaseQuery{
			Text: aws.String(query),
		},
		RetrievalConfiguration: &r.opts.RetrievalConfiguration,
	})
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, result := range page.RetrievalResults {
			docs = append(docs, schema.Document{
				PageContent: aws.ToString(result.Content.Text),
				Metadata: map[string]any{
					"location": aws.ToString(result.Location.S3Location.Uri),
					"score":    aws.ToFloat64(result.Score),
				},
			})
		}
	}

	return docs, nil
}

// Verbose returns the verbosity setting of the retriever.
func (r *BedrockKnowledgeBases) Verbose() bool {
	return r.opts.CallbackOptions.Verbose
}

// Callbacks returns the registered callbacks of the retriever.
func (r *BedrockKnowledgeBases) Callbacks() []schema.Callback {
	return r.opts.CallbackOptions.Callbacks
}
