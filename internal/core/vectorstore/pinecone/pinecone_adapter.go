package pinecone

import (
	"context"
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
)

type PineconeDB interface {
	ConvertToPineconeStruct(data []*model.CostAndUsage) []PineconeStruct
	Upsert(ctx context.Context, data []PineconeStruct) error
}
