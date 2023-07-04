package pinecone

import (
	"context"
	"crypto/tls"

	"github.com/google/uuid"
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

func (p *GRPCClient) Upsert(ctx context.Context, req *UpsertRequest) (*UpsertResponse, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", p.apiKey)

	pineconeVectors := make([]*pc.Vector, 0, len(req.Vectors))

	for i := 0; i < len(req.Vectors); i++ {
		metadataStruct, err := structpb.NewStruct(req.Vectors[i].Metadata)
		if err != nil {
			return nil, err
		}

		pineconeVectors = append(
			pineconeVectors,
			&pc.Vector{
				Id:       uuid.New().String(),
				Values:   float64ToFloat32(req.Vectors[i].Values),
				Metadata: metadataStruct,
			},
		)
	}

	pcRes, err := p.client.Upsert(ctx, &pc.UpsertRequest{
		Vectors:   pineconeVectors,
		Namespace: req.Namespace,
	})
	if err != nil {
		return nil, err
	}

	return &UpsertResponse{
		UpsertedCount: pcRes.UpsertedCount,
	}, nil
}

func (p *GRPCClient) Fetch(ctx context.Context, req *FetchRequest) (*FetchResponse, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", p.apiKey)

	pcRes, err := p.client.Fetch(ctx, &pc.FetchRequest{
		Ids:       req.IDs,
		Namespace: req.Namespace,
	})
	if err != nil {
		return nil, err
	}

	vectors := make(map[string]*Vector, len(pcRes.Vectors))
	for k, v := range pcRes.Vectors {
		vectors[k] = &Vector{
			ID:       v.Id,
			Values:   float32ToFloat64(v.Values),
			Metadata: v.Metadata.AsMap(),
		}
	}

	return &FetchResponse{
		Vectors:   nil,
		Namespace: pcRes.Namespace,
	}, nil
}

func (p *GRPCClient) Query(ctx context.Context, req *QueryRequest) (*QueryResponse, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", p.apiKey)

	filterStruct, err := structpb.NewStruct(req.Filter)
	if err != nil {
		return nil, err
	}

	pcRes, err := p.client.Query(ctx, &pc.QueryRequest{
		Namespace:       req.Namespace,
		TopK:            uint32(req.TopK),
		Filter:          filterStruct,
		IncludeValues:   req.IncludeValues,
		IncludeMetadata: req.IncludeMetadata,
		Queries: []*pc.QueryVector{{
			Values: float64ToFloat32(req.Vector),
		}},
	})
	if err != nil {
		return nil, err
	}

	matches := []*Match{}
	for _, m := range pcRes.Results[0].Matches {
		matches = append(matches, &Match{
			ID:       m.Id,
			Values:   float32ToFloat64(m.Values),
			Metadata: m.Metadata.AsMap(),
			Score:    float64(m.Score),
		})
	}

	return &QueryResponse{
		Namespace: pcRes.Results[0].Namespace,
		Matches:   matches,
	}, nil
}

func (p *GRPCClient) Close() error {
	return p.conn.Close()
}

func float64ToFloat32(input []float64) []float32 {
	output := make([]float32, len(input))
	for i, v := range input {
		output[i] = float32(v)
	}

	return output
}

func float32ToFloat64(input []float32) []float64 {
	result := make([]float64, len(input))
	for i, val := range input {
		result[i] = float64(val)
	}

	return result
}
