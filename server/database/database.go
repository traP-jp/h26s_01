package database

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Setup(mysqlConfig *mysql.Config) (*sqlx.DB, error) {
	if mysqlConfig == nil {
		return nil, errors.New("mysql config is nil")
	}

	db, err := sqlx.Connect("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("connecting mysql: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	if err := migrateTables(db.DB); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrating tables: %w", err)
	}

	return db, nil
}
