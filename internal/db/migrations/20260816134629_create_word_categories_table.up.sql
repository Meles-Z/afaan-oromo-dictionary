CREATE TABLE word_categories (
    word_id     INTEGER NOT NULL REFERENCES words(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (word_id, category_id)
);