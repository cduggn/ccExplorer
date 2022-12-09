package main

import (
	"github.com/cduggn/cloudcost/internal/pkg/cmd"
	"github.com/cduggn/cloudcost/internal/pkg/logger"
	"github.com/cduggn/cloudcost/internal/pkg/storage"
	_ "github.com/mattn/go-sqlite3"
)

var (
	database *storage.CostDataStorage
)

func main() {

	// create new database if instance does not already exist
	err := database.New("./cloudcost.db")
	if err != nil {
		logger.Error(err.Error())
	}

	cmd.Execute(database)
}
