package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func New(path string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	conn.SetMaxOpenConns(1) // SQLite: single-writer
	return conn, nil
}