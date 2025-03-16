package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {

	database, err := sql.Open("sqlite", "./main.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer database.Close()

	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	if err != nil {
		log.Fatalf("Error while preparing request: %v", err)
	}
	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		log.Fatalf("Error executing request: %v", err)
	}

	fmt.Println("The table was successfully created or already exists!")
}
