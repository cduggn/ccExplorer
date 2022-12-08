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

	record := storage.CostDataInsert{
		Dimension:   "test",
		Dimension2:  "test2",
		Tag:         "test3",
		MetricName:  "test4",
		Amount:      1.0,
		Unit:        "test5",
		Granularity: "test6",
		StartDate:   "test7",
		EndDate:     "test8",
	}

	storage.InsertCustomer(database, record)

}
