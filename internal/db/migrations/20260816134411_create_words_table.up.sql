CREATE TABLE words (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    english        TEXT NOT NULL,
    afaan_oromo    TEXT NOT NULL,
    part_of_speech TEXT,
    example_en     TEXT,
    example_om     TEXT,
    pronunciation  TEXT,
    created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_words_english ON words(english);
CREATE INDEX idx_words_afaan_oromo ON words(afaan_oromo);