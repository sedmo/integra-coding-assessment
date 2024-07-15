package db

import (
	"database/sql"
	"errors"
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

func InitDB(connector DBConnector) error {
	connectionString = os.Getenv("DATABASE_URL")
	if connectionString == "" {
		return errors.New("DATABASE_URL environment variable is not set")
	}

	var err error
	DB, err = connector.Open(connectionString)
	if err != nil {
		return err
	}

	if err := connector.Ping(DB); err != nil {
		return err
	}

	if err := connector.RunMigrationsDown(DB); err != nil {
		return err
	}

	if err := connector.RunMigrationsUp(DB); err != nil {
		return err
	}

	return nil
}
