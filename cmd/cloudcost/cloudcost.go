package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func main() {
	//cmd.Execute()
	createDB()
	insertCustomer()
}

func createDB() {

	if _, err := os.Stat("cloudcost.db"); os.IsNotExist(err) {
		db, err := os.Create("cloudcost.db")
		if err != nil {
			log.Fatal(err)
		}

		createTable()

		defer db.Close()

		log.Println("Database created")
	}

}

func createTable() {
	db, err := sql.Open("sqlite3", "./cloudcost.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
		CREATE TABLE customers (
		    till_id INTEGER PRIMARY KEY AUTOINCREMENT, 
		    client_id VARCHAR(64) NULL, 
		    first_name VARCHAR(255) NOT NULL, 
		    last_name VARCHAR(255) NOT NULL, 
		    guid VARCHAR(255) NULL, 
		    dob DATETIME NULL, type VARCHAR(1))
		    )
		    )
		`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s", err, sqlStmt)
		return
	}
}

func insertCustomer() {
	db, err := sql.Open("sqlite3", "./cloudcost.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO customers(client_id, first_name, last_name, guid, dob, type) values(?, ?, ?, ?, ?, ?)")
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
