package storage

import (
	"os"
)

var (
	//	createTableStmt = `
	//		CREATE TABLE cloudCostData (
	//		    till_id INTEGER PRIMARY KEY AUTOINCREMENT,
	//		    dimension VARCHAR(64) NULL,
	//		    dimension2 VARCHAR(255),
	//		    tag VARCHAR(255) NOT NULL,
	//		    metric_name VARCHAR(255) NOT NULL,
	//		    amount FLOAT NOT NULL,
	//		    unit VARCHAR(255) NOT NULL,
	//		    granularity VARCHAR(255) NOT NULL,
	//		    start_date DATETIME NOT NULL,
	//		    end_date DATETIME NOT NULL)
	//    `
	//	insertStmt = "INSERT INTO cloudCostData (dimension, dimension2, tag, metric_name, amount, unit, granularity, start_date, end_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	conn *CostDataStorage
)

func (c *CostDataStorage) New(dbName string) error {

	// if database already exists then no work required
	if conn != nil {
		return nil
	}

	// create the physical file if not already in place
	_, err := c.CreateFile(dbName)
	if err != nil {
		return DBError{msg: err.Error()}
	}

	// create the database name
	err = c.Set(dbName)
	if err != nil {
		return DBError{msg: err.Error()}
	}

	_, err = c.createCostDataTable()
	if err != nil {
		return DBError{msg: err.Error()}
	}

	return nil
}

// return 0 if creation was a success,
// return 1 if file already exists or return -1 if an error occured
func (c *CostDataStorage) CreateFile(dbName string) (int, error) {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		file, err := os.Create(dbName)
		if err != nil {
			return -1, &DBError{msg: "Could not create database"}
		}
		defer file.Close()
		return 0, nil
	}
	return 1, nil
}

// create a database with named provided as arg
func (c *CostDataStorage) Set(s string) error {
	//db, err := sql.Open("sqlite3", s)
	//if err != nil {
	//	return DBError{
	//		msg: err.Error(),
	//	}
	//}
	//c.SQLite = db
	return nil
}

// return -1 and or error if table was not created , return 0 if table was created
func (c *CostDataStorage) createCostDataTable() (int, error) {
	//_, err := c.SQLite.Exec(createTableStmt)
	//if err != nil {
	//	return -1, DBError{
	//		msg: err.Error(),
	//	}
	//}
	//logger.Info("Table created", zap.String("table", "cloudCostData"))
	return 0, nil
}

// return error and -1 if insert was not successful,
// return 0 if insert was successful
func (c *CostDataStorage) Insert(data CostDataInsert) (int, error) {

	//stmt, err := c.SQLite.Prepare(insertStmt)
	//if err != nil {
	//	logger.Error(err.Error())
	//	return -1, DBError{
	//		msg: err.Error(),
	//	}
	//}
	//defer stmt.Close()
	//
	//res, err := stmt.Exec(data.Dimension, data.Dimension2, data.Tag,
	//	data.MetricName, data.Amount, data.Unit, data.Granularity,
	//	data.StartDate, data.EndDate)
	//if err != nil {
	//	logger.Error(err.Error())
	//	return -1, DBError{
	//		msg: err.Error(),
	//	}
	//}
	//_, err = res.LastInsertId()
	//if err != nil {
	//	logger.Error(err.Error())
	//	return -1, DBError{
	//		msg: err.Error(),
	//	}
	//}
	return 0, nil
}
