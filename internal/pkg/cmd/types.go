package cmd

import (
	"github.com/cduggn/cloudcost/internal/pkg/billing"
	"github.com/cduggn/cloudcost/internal/pkg/storage"
)

type DB struct {
	*storage.CostDataStorage
}

var report *billing.CostAndUsageReport
