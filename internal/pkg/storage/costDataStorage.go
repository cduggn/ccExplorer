package storage

import (
	"database/sql"
	"fmt"
	"github.com/cduggn/cloudcost/internal/pkg/logger"
	"go.uber.org/zap"
	"log"
	"os"
)

var (
	dbType          = "sqlite3"
	dbName          = "cloudcost.db"
	dbLocation      = "./"
	createTableStmt = `
		CREATE TABLE cloudCostData (
		    till_id INTEGER PRIMARY KEY AUTOINCREMENT, 
		    dimension VARCHAR(64) NULL, 
		    dimension2 VARCHAR(255),
		    tag VARCHAR(255) NOT NULL, 
		    metric_name VARCHAR(255) NOT NULL,
		    amount FLOAT NOT NULL,
		    unit VARCHAR(255) NOT NULL,
		    granularity VARCHAR(255) NOT NULL,
		    start_date DATETIME NOT NULL,
		    end_date DATETIME NOT NULL)
    `
	insertStmt = "INSERT INTO cloudCostData (dimension, dimension2, tag, metric_name, amount, unit, granularity, start_date, end_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
)

func (cd CostDataStorage) NewPersistentStorage() {

	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		db, err := os.Create(dbName)
		if err != nil {
			logger.Error(err.Error())
		}
		defer db.Close()

		logger.Info("Database created", zap.String("db", dbName))
	}
}

func (cd CostDataStorage) CreateCostDataTable() {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s%s", dbLocation, dbName))
	if err != nil {
		logger.Error(err.Error())
	}
	defer db.Close()

	_, err = db.Exec(createTableStmt)
	if err != nil {
		log.Printf("%q: %s", err, createTableStmt)
		return
	}
	logger.Info("Table created", zap.String("table", "cloudCostData"))
}

func (cd CostDataStorage) InsertCustomer() {
	db, err := sql.Open(dbType, fmt.Sprintf("%s%s", dbLocation, dbName))
	if err != nil {
		logger.Error(err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare(insertStmt)
	if err != nil {
		logger.Error(err.Error())
	}
	defer stmt.Close()

	res, err := stmt.Exec("dimension", "dimension2", "tag", "metric_name", 1.0, "unit", "granularity", "start_date", "end_date")
	if err != nil {
		logger.Error(err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Info("Row added", zap.Int64("rowId", id))
}
