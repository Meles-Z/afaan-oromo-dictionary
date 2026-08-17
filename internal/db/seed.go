package db

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed seed/words.json
var seedFS embed.FS

type seedWord struct {
	English    string `json:"english"`
	AfaanOromo string `json:"afaanOromo"`
}

// SeedIfEmpty populates the words table from embedded seed data, but only
// if the table is currently empty — safe to call on every startup, and
// never overwrites user-added or user-edited entries.
func SeedIfEmpty(conn *sql.DB) error {
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM words`).Scan(&count); err != nil {
		return fmt.Errorf("count words: %w", err)
	}
	if count > 0 {
		return nil // already seeded (or user has their own data) - do nothing
	}

	data, err := seedFS.ReadFile("seed/words.json")
	if err != nil {
		return fmt.Errorf("read embedded seed data: %w", err)
	}

	var words []seedWord
	if err := json.Unmarshal(data, &words); err != nil {
		return fmt.Errorf("parse seed data: %w", err)
	}

	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("begin seed transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO words (english, afaan_oromo) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("prepare seed insert: %w", err)
	}
	defer stmt.Close()

	for _, w := range words {
		if _, err := stmt.Exec(w.English, w.AfaanOromo); err != nil {
			return fmt.Errorf("insert seed word %q: %w", w.English, err)
		}
	}

	return tx.Commit()
}
