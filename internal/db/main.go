package db

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// DBConnection is an interface that abstracts the database operations.
type DBConnection interface {
	Open(driverName, dataSourceName string) (*sql.DB, error)
	Close(db *sql.DB) error
}

// SQLDBConnection is the implementation for SQL database operations.
type SQLDBConnection struct{}

// Open establishes a new database connection using the specified driver and data source name.
func (s *SQLDBConnection) Open(driverName, dataSourceName string) (*sql.DB, error) {
	return sql.Open(driverName, dataSourceName)
}

// Close terminates the given database connection.
func (s *SQLDBConnection) Close(db *sql.DB) error {
	return db.Close()
}

var (
	db                    *sql.DB
	dbMu                  sync.Mutex
	dbConn                DBConnection = &SQLDBConnection{}
	defaultDriverName                  = "sqlite3"
	defaultDataSourceName              = "database.db"
)

const (
	errDBOpen               = "error opening database connection"
	errDBAlreadyInitialized = "database already initialized"
	errDBClose              = "error closing database"
	errDBNotInitialized     = "database not initialized"
)

// InitDB initializes the database connection. It takes optional parameters for the driver name and data source name.
func InitDB(driverName, dataSourceName string) error {
	dbMu.Lock()
	defer dbMu.Unlock()

	if db != nil {
		return errors.New(errDBAlreadyInitialized)
	}

	if driverName == "" {
		driverName = defaultDriverName
	}
	if dataSourceName == "" {
		dataSourceName = defaultDataSourceName
	}

	var err error
	db, err = dbConn.Open(driverName, dataSourceName)
	if err != nil {
		return fmt.Errorf("%s: %w", errDBOpen, err)
	}

	return nil
}

// CloseDB closes the database connection if it is initialized.
func CloseDB() error {
	dbMu.Lock()
	defer dbMu.Unlock()

	if db == nil {
		return errors.New(errDBNotInitialized)
	}

	if err := dbConn.Close(db); err != nil {
		return fmt.Errorf("%s: %w", errDBClose, err)
	}

	db = nil
	return nil
}

// GetDB returns the database connection if it is initialized.
func GetDB() (*sql.DB, error) {
	dbMu.Lock()
	defer dbMu.Unlock()

	if db == nil {
		return nil, errors.New(errDBNotInitialized)
	}

	return db, nil
}

// SetDBConnection allows setting a different DBConnection implementation for testing.
func SetDBConnection(conn DBConnection) {
	dbMu.Lock()
	defer dbMu.Unlock()
	dbConn = conn
}
