package db

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	dbMu sync.Mutex
)

const (
	defaultDBFilePath = "verbs.db"

	errDBOpen               = "error opening database connection"
	errDBAlreadyInitialized = "database already initialized"
	errDBClose              = "error closing database"
	errDBNotInitialized     = "database not initialized"
)

// InitDB initializes the database connection.
func InitDB(filePath ...string) error {
	dbMu.Lock()
	defer dbMu.Unlock()

	if db != nil {
		return errors.New(errDBAlreadyInitialized)
	}

	// Use default file path if none is provided
	dbFilePath := defaultDBFilePath
	if len(filePath) > 0 {
		dbFilePath = filePath[0]
	}

	var err error
	db, err = sql.Open("sqlite3", dbFilePath)
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

	if err := db.Close(); err != nil {
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
