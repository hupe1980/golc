package pinecone

type Vector struct {
	ID       string         `json:"id"`
	Values   []float32      `json:"values"`
	Metadata map[string]any `json:"metadata"`
}

// UpsertRequest represents the parameters for an upsert vectors request.
// See https://docs.pinecone.io/reference/upsert for more information.
type UpsertRequest struct {
	Vectors   []*Vector `json:"vectors"`
	Namespace string    `json:"namespace"`
}

// UpsertResponse represents the response from an upsert vectors request.
type UpsertResponse struct {
	UpsertedCount int `json:"upsertedCount"`
}
