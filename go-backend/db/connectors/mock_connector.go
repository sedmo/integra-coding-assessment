package connectors

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

type MockConnector struct {
	DB      *sql.DB
	Sqlmock sqlmock.Sqlmock
}

func (m *MockConnector) Open(connectionString string) (*sql.DB, error) {
	var err error
	m.DB, m.Sqlmock, err = sqlmock.New()
	return m.DB, err
}

func (m *MockConnector) Ping(db *sql.DB) error {
	return db.Ping()
}

func (m *MockConnector) RunMigrationsUp(db *sql.DB) error {
	return nil
}

func (m *MockConnector) RunMigrationsDown(db *sql.DB) error {
	return nil
}
