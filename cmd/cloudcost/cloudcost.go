package main

import (
	"github.com/cduggn/cloudcost/internal/pkg/cmd"
	"github.com/cduggn/cloudcost/internal/pkg/storage"
	_ "github.com/mattn/go-sqlite3"
)

var (
	database *storage.CostDataStorage
)

func main() {
	cmd.Execute()
	database = storage.New("sqlite3", "./cloudcost.db")
	storage.InsertCustomer(database)

}
