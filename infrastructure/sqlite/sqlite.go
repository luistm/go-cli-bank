package sqlite

import (
	"database/sql"
	"errors"
	"os"

	"github.com/luistm/go-bank-cli/infrastructure"
	"github.com/luistm/go-bank-cli/lib/categories"
	// To init the database driver
	_ "github.com/mattn/go-sqlite3"
)

var errConnectionIsNil = errors.New("sqlite database is <nil>")
var errInvalidConfiguration = errors.New("sqlite configuration parameters are invalid")
var errFailedToCreatedDB = errors.New("failed to create database")

// New creates a new instance of sqlite
func New(path string, name string) (infrastructure.Storage, error) {

	if name == "" {
		return &Sqlite{}, errInvalidConfiguration
	}

	if path == "" {
		return &Sqlite{}, errInvalidConfiguration
	}

	if err := validatePath(path); err != nil {
		return &Sqlite{}, errInvalidConfiguration
	}

	db, err := sql.Open("sqlite3", path+"/"+name+".db")
	if err != nil {
		return nil, err
	}

	// Create table in order to create the database file
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS categories
	(id INTEGER NOT NULL PRIMARY KEY, name TEXT);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return &Sqlite{}, errFailedToCreatedDB
	}

	s := &Sqlite{db}
	return s, nil
}

// validatePath validates that the database directory or file exists and it is
// in a proper format to be used
func validatePath(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	} else {
		if err != nil {
			return err
		}
		fileInfo, _ = os.Stat(path)
		if !fileInfo.Mode().IsDir() {
			return err
		}
	}

	return nil
}

// Sqlite contains a sqlite database
type Sqlite struct {
	db *sql.DB
}

// Close closes the connection with the sqlite database
func (s *Sqlite) Close() error {
	if s.db == nil {
		return errConnectionIsNil
	}

	return s.db.Close()
}

// Execute is to execute an sql statement
func (s *Sqlite) Execute(statement string, values ...interface{}) error {
	if s.db == nil {
		return errConnectionIsNil
	}

	tx, _ := s.db.Begin()
	_, err := tx.Exec(statement, values...)
	if err != nil {
		// TODO: tx.Rollback
		return &ErrDataBase{err.Error()}
	}
	tx.Commit()

	return nil
}

// Query fetches data from the database
func (s *Sqlite) Query(statement string) (categories.IRow, error) {
	if s.db == nil {
		return nil, errConnectionIsNil
	}

	rows, err := s.db.Query(statement)
	if err != nil {
		return nil, &ErrDataBase{err.Error()}
	}

	return rows, nil
}

// LEGACY ---------------------------------------------------------------------

// ErrDataBase to be used when the infrastructure
// database returns error
type ErrDataBase struct {
	s string
}

func (e *ErrDataBase) Error() string {
	return e.s
}