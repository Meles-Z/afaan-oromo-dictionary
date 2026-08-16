package dictionary

import (
	"database/sql"
	"fmt"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Search(query, direction string) ([]Word, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []Word{}, nil
	}

	column := "english"
	if direction == "om" {
		column = "afaan_oromo"
	}

	ftsQuery := fmt.Sprintf(`"%s"*`, strings.ReplaceAll(query, `"`, `""`))

	rows, err := r.db.Query(fmt.Sprintf(`
    SELECT w.id, w.english, w.afaan_oromo, w.part_of_speech, w.example_en, w.example_om,
           w.pronunciation, w.created_at, w.updated_at
    FROM words_fts f
    JOIN words w ON w.id = f.rowid
    WHERE f.%s MATCH ?
    ORDER BY
      CASE WHEN w.%s = ?2 COLLATE NOCASE THEN 0 ELSE 1 END,
      f.rank
    LIMIT 50
`, column, column), ftsQuery, query)
	if err != nil {
		return nil, fmt.Errorf("search words: %w", err)
	}
	defer rows.Close()

	words, err := scanWords(rows)
	if err != nil {
		return nil, err
	}

	// Fallback: no FTS matches (e.g. typo, or a substring mid-word) — try a
	// plain contains search so the user isn't left with a false "not found"
	if len(words) == 0 {
		return r.searchContains(query, column)
	}

	return words, nil
}

func (r *Repository) GetByID(id int) (*Word, error) {
	var w Word
	var partOfSpeech, exampleEn, exampleOm, pronunciation sql.NullString

	err := r.db.QueryRow(`
		SELECT id, english, afaan_oromo, part_of_speech, example_en, example_om,
		       pronunciation, created_at, updated_at
		FROM words WHERE id = ?
	`, id).Scan(&w.ID, &w.English, &w.AfaanOromo, &partOfSpeech,
		&exampleEn, &exampleOm, &pronunciation, &w.CreatedAt, &w.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get word by id: %w", err)
	}

	w.PartOfSpeech = partOfSpeech.String
	w.ExampleEn = exampleEn.String
	w.ExampleOm = exampleOm.String
	w.Pronunciation = pronunciation.String

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

func (r *Repository) searchContains(query, column string) ([]Word, error) {
	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT id, english, afaan_oromo, part_of_speech, example_en, example_om,
		       pronunciation, created_at, updated_at
		FROM words
		WHERE %s LIKE ?
		ORDER BY %s
		LIMIT 50
	`, column, column), "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("fallback search: %w", err)
	}
	defer rows.Close()
	return scanWords(rows)
}

func scanWords(rows *sql.Rows) ([]Word, error) {
	var words []Word
	for rows.Next() {
		var w Word
		var partOfSpeech, exampleEn, exampleOm, pronunciation sql.NullString

		if err := rows.Scan(&w.ID, &w.English, &w.AfaanOromo, &partOfSpeech,
			&exampleEn, &exampleOm, &pronunciation, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan word: %w", err)
		}
		w.PartOfSpeech = partOfSpeech.String
		w.ExampleEn = exampleEn.String
		w.ExampleOm = exampleOm.String
		w.Pronunciation = pronunciation.String
		words = append(words, w)
	}
	return words, rows.Err()
}
