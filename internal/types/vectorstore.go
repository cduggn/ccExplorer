package types

type VectorStoreInput struct {
	Items []VectorStoreItem
}

type VectorStoreItem struct {
	ID              string
	EmbeddingText   string
	EmbeddingVector []float32
	Metadata        VectorStoreItemMetadata
}

type VectorStoreItemMetadata struct {
	StartDate   string
	EndDate     string
	Granularity string
	Dimensions  string
	Tags        string
	Cost        string
}

type UpsertResponse struct {
	UpsertedCount int `json:"upsertedCount"`
}
