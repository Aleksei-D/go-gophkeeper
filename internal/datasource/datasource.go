package datasource

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const (
	Server = iota
	Agent
)

func NewDatabase(databaseURI string, typeDB int) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURI)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	var migrations func(db *sql.DB) error
	switch typeDB {
	case Server:
		migrations = upServerMigrations
	case Agent:
		migrations = upAgentMigrations
	}

	if err := migrations(db); err != nil {
		return nil, err
	}
	return db, nil
}
