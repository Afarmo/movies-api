package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := " CREATE TABLE IF NOT EXISTS users ( id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,name TEXT    )"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to SQLite!")
}
