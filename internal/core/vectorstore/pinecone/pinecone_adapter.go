package pinecone

import (
	"context"
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
)

type PineconeDB interface {
	ConvertToVectorStoreItem(r model.CostAndUsageOutputType) []model.
		VectorStoreItem
	Upsert(ctx context.Context, data []PineconeStruct) error
}
