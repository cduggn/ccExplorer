package cmd

import (
	"github.com/cduggn/cloudcost/internal/pkg/storage"
)

type DB struct {
	*storage.CostDataStorage
}
