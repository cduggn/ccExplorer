package pinecone

import (
	"context"
	"github.com/cduggn/ccexplorer/internal/types"
)

type PineconeDB interface {
	ConvertToVectorStoreItem(r types.CostAndUsageOutputType) []*types.VectorStoreItem
	Upsert(ctx context.Context, data []PineconeStruct) error
}
