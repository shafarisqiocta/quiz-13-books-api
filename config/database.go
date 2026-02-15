package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	connStr := "host=localhost port=5432 user=postgres password=test123 dbname=books_api sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database not responding:", err)
	}

	fmt.Println("Database connected!")

	DB = db
}
