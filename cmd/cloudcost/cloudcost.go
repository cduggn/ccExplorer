package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
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
		    client_id VARCHAR(64) NULL, 
		    first_name VARCHAR(255) NOT NULL, 
		    last_name VARCHAR(255) NOT NULL, 
		    guid VARCHAR(255) NULL, 
		    dob DATETIME NULL, type VARCHAR(1))
		    )
		    )
    `
	insertStmt = "INSERT INTO cloudCostData(client_id, first_name, last_name, guid, dob, type) values(?, ?, ?, ?, ?, ?)"
)

func main() {
	//cmd.Execute()
	createDatabase()
	insertCustomer()
}

func createDatabase() {

	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		db, err := os.Create(dbName)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		createCostDataTable()
		log.Println("Database created")
	}
}

func createCostDataTable() {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s%s", dbLocation, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(createTableStmt)
	if err != nil {
		log.Printf("%q: %s", err, createTableStmt)
		return
	}
}

func insertCustomer() {
	db, err := sql.Open(dbType, fmt.Sprintf("%s%s", dbLocation, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare(insertStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec("1", "John", "Doe", "123456789", "1980-01-01", "C")
	if err != nil {
		log.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}
