package database

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

func migrateTables(db *sql.DB) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("mysql"); err != nil {
		return fmt.Errorf("setting dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("applying migrations: %w", err)
	}

	return nil
}
