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
