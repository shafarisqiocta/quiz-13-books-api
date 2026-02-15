package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Fallback ke localhost untuk development
	if host == "" {
		host = "localhost"
		port = "5432"
		user = "postgres"
		password = "test123"
		dbname = "books_api"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database not responding:", err)
	}

	fmt.Println(" Database connected!")

	DB = db
}
