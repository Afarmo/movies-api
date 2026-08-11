package main

import (
	"log"
	"movies-api/internal/database"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := database.NewDataBase("movieapi.db")
	if err != nil {
		log.Fatal(err)
	}
	if err := database.InitializeQuery(db); err != nil {
		log.Fatal(err)
	}
}
