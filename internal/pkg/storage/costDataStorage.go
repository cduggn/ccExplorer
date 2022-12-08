package storage

import (
	"database/sql"
	"github.com/cduggn/cloudcost/internal/pkg/logger"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"log"
	"os"
)

var (
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
	database   *CostDataStorage
)

func newPersistentStorage(dbName string) int {
	file, err := os.Create(dbName)
	if err != nil {
		logger.Error(err.Error())
		return -1
	}
	defer file.Close()
	logger.Info("Database created", zap.String("db", dbName))
	return 0
}

func newConnection(dbType string, dbName string) (*CostDataStorage, error) {
	db, err := sql.Open(dbType, dbName)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	database = &CostDataStorage{db}

	return database, nil
}

func New(driverName, dbName string) *CostDataStorage {

	if database != nil {
		return database
	}

	newDatabase := false
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		newDatabase = true
		res := newPersistentStorage(dbName)
		if res == -1 {
			logger.Error("Could not create database")
			return nil
		}
	}

	db, err := newConnection(driverName, dbName)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	if newDatabase {
		res := createCostDataTable(db)
		if res == -1 {
			logger.Error("Could not create table")
			return nil
		}
	}
	return db
}

func createCostDataTable(db *CostDataStorage) int {
	_, err := db.Exec(createTableStmt)
	if err != nil {
		log.Printf("%q: %s", err, createTableStmt)
		return -1
	}
	logger.Info("Table created", zap.String("table", "cloudCostData"))
	return 0
}

func InsertCustomer(db *CostDataStorage, data CostDataInsert) int {

	stmt, err := db.Prepare(insertStmt)
	if err != nil {
		logger.Error(err.Error())
		return -1
	}
	defer stmt.Close()

	res, err := stmt.Exec(data.Dimension, data.Dimension2, data.Tag, data.MetricName, data.Amount, data.Unit, data.Granularity, data.StartDate, data.EndDate)
	if err != nil {
		logger.Error(err.Error())
		return -1
	}
	id, err := res.LastInsertId()
	if err != nil {
		logger.Error(err.Error())
		return -1
	}
	logger.Info("Row added", zap.Int64("rowId", id))
	return 0
}



