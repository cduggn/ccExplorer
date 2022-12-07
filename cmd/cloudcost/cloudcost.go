package main

import (
	"fmt"
	"github.com/cduggn/cloudcost/internal/pkg/storage"
	_ "github.com/mattn/go-sqlite3"
)

var (
	costDataStorage storage.CostDataStorage
)

func main() {
	//cmd.Execute()
	createStorageBackend(costDataStorage)
	insertCustomer()
	fmt.Println("Work done!")
}

func createStorageBackend(st storage.Persistent) {
	st.NewPersistentStorage()
	st.CreateCostDataTable()
}

func insertCustomer() {
	costDataStorage.InsertCustomer()
}
