package connectors

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Import for file source
	_ "github.com/lib/pq"                                // Import for Postgres driver
)

type PostgresConnector struct{}

func (p *PostgresConnector) Open(connectionString string) (*sql.DB, error) {
	return sql.Open("postgres", connectionString)
}

func (p *PostgresConnector) Ping(db *sql.DB) error {
	return db.Ping()
}

func (p *PostgresConnector) RunMigrationsUp(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file:///app/db/migrations", "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	log.Println("Database migrated successfully")
	return nil
}

func (p *PostgresConnector) RunMigrationsDown(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file:///app/db/migrations", "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	log.Println("Database rolled back successfully")
	return nil
}
