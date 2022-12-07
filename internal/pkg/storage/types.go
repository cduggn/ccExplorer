package storage

type Persistent interface {
	NewPersistentStorage()
	CreateCostDataTable()
	InsertCustomer()
}

type CostDataStorage struct {
}
