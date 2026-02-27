package database

import (
	"database/sql"
	"fmt"
	"os"
)

type DatabasePostgres struct{}

func NewDatabasePostgres() *DatabasePostgres {
	return &DatabasePostgres{}
}

func (d DatabasePostgres) Connect() (*sql.DB, error) {
	// Retrieve credentials from environment variables
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := "localhost" //os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	fmt.Printf("Port: %s\n", port)

	// TODO make configurable
	sslmode := "disable" // local dev
	timezone := "America/New_York"

	dataSourceName := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s TimeZone=%s",
		user,
		pass,
		dbName,
		host,
		port,
		sslmode,
		timezone,
	)

	db, err := sql.Open("postgres", dataSourceName)
	return db, err
}
