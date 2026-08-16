CREATE VIRTUAL TABLE IF NOT EXISTS words_fts USING fts5(
    english,
    afaan_oromo,
    content='words',
    content_rowid='id',
    tokenize='unicode61'
);

INSERT INTO words_fts(rowid, english, afaan_oromo)
SELECT id, english, afaan_oromo FROM words;

CREATE TRIGGER words_ai AFTER INSERT ON words BEGIN
  INSERT INTO words_fts(rowid, english, afaan_oromo) VALUES (new.id, new.english, new.afaan_oromo);
END;

CREATE TRIGGER words_ad AFTER DELETE ON words BEGIN
  INSERT INTO words_fts(words_fts, rowid, english, afaan_oromo) VALUES('delete', old.id, old.english, old.afaan_oromo);
END;

CREATE TRIGGER words_au AFTER UPDATE ON words BEGIN
  INSERT INTO words_fts(words_fts, rowid, english, afaan_oromo) VALUES('delete', old.id, old.english, old.afaan_oromo);
  INSERT INTO words_fts(rowid, english, afaan_oromo) VALUES (new.id, new.english, new.afaan_oromo);
END;