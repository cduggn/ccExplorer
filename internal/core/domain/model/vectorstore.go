package model

type VectorStoreInput struct {
	Items []VectorStoreItem
}

type VectorStoreItem struct {
	EmbeddingText   string
	EmbeddingVector []float32
	Metadata        VectorStoreItemMetadata
}

type VectorStoreItemMetadata struct {
	StartDate   string
	Granularity string
	Dimensions  string
	Tags        string
}
