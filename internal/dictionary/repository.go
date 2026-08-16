package dictionary

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Search does a prefix match on english. Direction ("en" or "om") decides
// which column to search against.
func (r *Repository) Search(query, direction string) ([]Word, error) {
	column := "english"
	if direction == "om" {
		column = "afaan_oromo"
	}

	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT id, english, afaan_oromo, part_of_speech, example_en, example_om,
		       pronunciation, created_at, updated_at
		FROM words
		WHERE %s LIKE ?
		ORDER BY %s
		LIMIT 50
	`, column, column), query+"%")
	if err != nil {
		return nil, fmt.Errorf("search words: %w", err)
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var w Word
		if err := rows.Scan(&w.ID, &w.English, &w.AfaanOromo, &w.PartOfSpeech,
			&w.ExampleEn, &w.ExampleOm, &w.Pronunciation, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan word: %w", err)
		}
		words = append(words, w)
	}
	return words, rows.Err()
}

func (r *Repository) GetByID(id int) (*Word, error) {
	var w Word
	err := r.db.QueryRow(`
		SELECT id, english, afaan_oromo, part_of_speech, example_en, example_om,
		       pronunciation, created_at, updated_at
		FROM words WHERE id = ?
	`, id).Scan(&w.ID, &w.English, &w.AfaanOromo, &w.PartOfSpeech,
		&w.ExampleEn, &w.ExampleOm, &w.Pronunciation, &w.CreatedAt, &w.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get word by id: %w", err)
	}
	return &w, nil
}

func (r *Repository) Create(w WordInput) (int, error) {
	res, err := r.db.Exec(`
		INSERT INTO words (english, afaan_oromo, part_of_speech, example_en, example_om, pronunciation)
		VALUES (?, ?, ?, ?, ?, ?)
	`, w.English, w.AfaanOromo, w.PartOfSpeech, w.ExampleEn, w.ExampleOm, w.Pronunciation)
	if err != nil {
		return 0, fmt.Errorf("insert word: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id: %w", err)
	}

	if err := r.setCategories(int(id), w.CategoryIDs); err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *Repository) Update(id int, w WordInput) error {
	_, err := r.db.Exec(`
		UPDATE words
		SET english = ?, afaan_oromo = ?, part_of_speech = ?, example_en = ?,
		    example_om = ?, pronunciation = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, w.English, w.AfaanOromo, w.PartOfSpeech, w.ExampleEn, w.ExampleOm, w.Pronunciation, id)
	if err != nil {
		return fmt.Errorf("update word: %w", err)
	}
	return r.setCategories(id, w.CategoryIDs)
}

func (r *Repository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM words WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete word: %w", err)
	}
	return nil
}

// setCategories replaces all category links for a word (delete + reinsert —
// simplest correct approach for a small join table).
func (r *Repository) setCategories(wordID int, categoryIDs []int) error {
	if _, err := r.db.Exec(`DELETE FROM word_categories WHERE word_id = ?`, wordID); err != nil {
		return fmt.Errorf("clear categories: %w", err)
	}
	for _, catID := range categoryIDs {
		if _, err := r.db.Exec(
			`INSERT INTO word_categories (word_id, category_id) VALUES (?, ?)`,
			wordID, catID,
		); err != nil {
			return fmt.Errorf("link category: %w", err)
		}
	}
	return nil
}

func (r *Repository) ListCategories() ([]Category, error) {
	rows, err := r.db.Query(`SELECT id, name FROM categories ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	var cats []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}