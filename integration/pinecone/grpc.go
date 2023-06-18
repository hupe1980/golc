package pinecone

import (
	"context"
	"crypto/tls"

	"github.com/google/uuid"
	"github.com/hupe1980/golc/util"
	pc "github.com/pinecone-io/go-pinecone/pinecone_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
)

type GRPCClient struct {
	apiKey string
	conn   *grpc.ClientConn
	client pc.VectorServiceClient
}

func NewGRPCClient(apiKey string, endpoint Endpoint) (*GRPCClient, error) {
	target := endpoint.String()

	conn, err := grpc.Dial(
		target,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
		})),
		grpc.WithAuthority(target),
		grpc.WithBlock(),
	)

	if err != nil {
		return nil, err
	}

	return &GRPCClient{
		apiKey: apiKey,
		conn:   conn,
		client: pc.NewVectorServiceClient(conn),
	}, nil
}

func (p *GRPCClient) Upsert(ctx context.Context, req *pc.UpsertRequest) (*pc.UpsertResponse, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", p.apiKey)
	return p.client.Upsert(ctx, req)
}

func (p *GRPCClient) Fetch(ctx context.Context, req *pc.FetchRequest) (*pc.FetchResponse, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", p.apiKey)
	return p.client.Fetch(ctx, req)
}

func (p *GRPCClient) Query(ctx context.Context, req *pc.QueryRequest) (*pc.QueryResponse, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", p.apiKey)
	return p.client.Query(ctx, req)
}

func (p *GRPCClient) Close() error {
	return p.conn.Close()
}

func ToPineconeGRPCVectors(vectors [][]float64, metadata []map[string]any) ([]*pc.Vector, error) {
	pineconeVectors := make([]*pc.Vector, 0, len(vectors))

	for i := 0; i < len(vectors); i++ {
		metadataStruct, err := structpb.NewStruct(metadata[i])
		if err != nil {
			return nil, err
		}

		pineconeVectors = append(
			pineconeVectors,
			&pc.Vector{
				Id: uuid.New().String(),
				Values: util.Map(vectors[i], func(v float64, _ int) float32 {
					return float32(v)
				}),
				Metadata: metadataStruct,
			},
		)
	}

	return pineconeVectors, nil
}
