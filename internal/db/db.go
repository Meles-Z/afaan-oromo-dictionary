package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// New opens the app's SQLite database, storing it in the OS-appropriate
// user data directory so it works correctly whether run in dev or installed.
func New(filename string) (*sql.DB, error) {
	dir, err := appDataDir()
	if err != nil {
		return nil, fmt.Errorf("resolve app data dir: %w", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create app data dir: %w", err)
	}

	path := filepath.Join(dir, filename)

	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	conn.SetMaxOpenConns(1)

	return conn, nil
}

func appDataDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "afaan-oromo-dictionary"), nil
}
