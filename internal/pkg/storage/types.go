package storage

import "database/sql"

type Persistent interface {
	NewPersistentStorage()
	CreateCostDataTable()
	InsertCustomer()
}

type CostDataStorage struct {
	*sql.DB
}

type CostDataInsert struct {
	Dimension   string
	Dimension2  string
	Tag         string
	MetricName  string
	Amount      float64
	Unit        string
	Granularity string
	StartDate   string
	EndDate     string
}
