package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/Masterminds/squirrel"
)

var DB *sql.DB
var Psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

var connectionString string

type DBConnector interface {
	Open(connectionString string) (*sql.DB, error)
	Ping(db *sql.DB) error
	RunMigrationsUp(db *sql.DB) error
	RunMigrationsDown(db *sql.DB) error
}

func InitDB(connector DBConnector) {
	connectionString = os.Getenv("DATABASE_URL")
	if connectionString == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	var err error
	DB, err = connector.Open(connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := connector.Ping(DB); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	if err := connector.RunMigrationsDown(DB); err != nil {
		log.Fatalf("Error running migrations down: %v", err)
	}

	if err := connector.RunMigrationsUp(DB); err != nil {
		log.Fatalf("Error running migrations up: %v", err)
	}
}
