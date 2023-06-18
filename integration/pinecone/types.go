package pinecone

type Vector struct {
	ID       string         `json:"id"`
	Values   []float32      `json:"values"`
	Metadata map[string]any `json:"metadata"`
}

// UpsertRequest represents the parameters for an upsert vectors request.
// See https://docs.pinecone.io/reference/upsert for more informations.
type UpsertRequest struct {
	Vectors   []*Vector `json:"vectors"`
	Namespace string    `json:"namespace"`
}

// UpsertResponse represents the response from an upsert vectors request.
type UpsertResponse struct {
	UpsertedCount uint32 `json:"upsertedCount"`
}

// FetchRequest represents the parameters for a fetch vectors request.
// https://docs.pinecone.io/reference/fetch for more informations.
type FetchRequest struct {
	IDs       []string `json:"ids"`
	Namespace string   `json:"namespace"`
}

// FetchResponse represents the response from a fetch vectors request.
type FetchResponse struct {
	Vectors   map[string]*Vector `json:"vectors"`
	Namespace string             `json:"namespace"`
}

// QueryRequest represents the parameters for a query request.
// See https://docs.pinecone.io/reference/query for more information.
type QueryRequest struct {
	Filter          map[string]any `json:"filter"`
	IncludeValues   bool           `json:"includeValues"`
	IncludeMetadata bool           `json:"includeMetadata"`
	Vector          []float32      `json:"vector"`
	Namespace       string         `json:"namespace"`
	TopK            int64          `json:"topK"`
	ID              string         `json:"id"`
}

type Match struct {
	ID       string         `json:"id"`
	Values   []float32      `json:"values"`
	Metadata map[string]any `json:"metadata"`
	Score    float32        `json:"score"`
}

// QueryResponse represents the response from a query request.
type QueryResponse struct {
	Matches   []*Match `json:"matches"`
	Namespace string   `json:"namespace"`
}
