package storage

import "database/sql"

type Persistent interface {
	NewPersistentStorage()
	CreateCostDataTable()
	InsertCustomer()
}

type CostDataStorage struct {
	SQLite *sql.DB
}

type DBError struct {
	msg string
}

func (e DBError) Error() string {
	return e.msg
}

type CostDataInsert struct {
	Dimension   string
	Dimension2  string
	Tag         string
	MetricName  string
	Amount      string
	Unit        string
	Granularity string
	StartDate   string
	EndDate     string
}
