package pinecone

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hupe1980/golc/util"
)

type Endpoint struct {
	IndexName   string
	ProjectName string
	Environment string
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s-%s.svc.%s.pinecone.io:443", e.IndexName, e.ProjectName, e.Environment)
}

type Options struct {
	UseGRPC bool
}

type Client interface {
	Upsert(ctx context.Context, req *UpsertRequest) (*UpsertResponse, error)
	Fetch(ctx context.Context, req *FetchRequest) (*FetchResponse, error)
	Query(ctx context.Context, req *QueryRequest) (*QueryResponse, error)
	Close() error
}

func New(apiKey string, endpoint Endpoint, optFns ...func(o *Options)) (Client, error) {
	opts := Options{
		UseGRPC: false,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.UseGRPC {
		return NewGRPCClient(apiKey, endpoint)
	}

	return NewRestClient(apiKey, endpoint)
}

func ToPineconeVectors(vectors [][]float64, metadata []map[string]any) ([]*Vector, error) {
	pineconeVectors := make([]*Vector, 0, len(vectors))

	for i := 0; i < len(vectors); i++ {
		pineconeVectors = append(
			pineconeVectors,
			&Vector{
				ID: uuid.New().String(),
				Values: util.Map(vectors[i], func(v float64, _ int) float32 {
					return float32(v)
				}),
				Metadata: metadata[i],
			},
		)
	}

	return pineconeVectors, nil
}
