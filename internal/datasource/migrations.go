package datasource

import (
	"database/sql"
	"embed"
	"io/fs"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/server/*.sql
var embedServerMigrations embed.FS

//go:embed migrations/agent/*.sql
var embedAgentMigrations embed.FS

func upServerMigrations(db *sql.DB) error {
	return upMigrations(db, embedServerMigrations)
}

func upAgentMigrations(db *sql.DB) error {
	return upMigrations(db, embedAgentMigrations)
}

func upMigrations(db *sql.DB, fsys fs.FS) error {
	goose.SetBaseFS(fsys)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	return nil
}
