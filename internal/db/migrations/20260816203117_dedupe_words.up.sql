
DELETE FROM words
WHERE id NOT IN (
  SELECT MIN(id)
  FROM words
  GROUP BY LOWER(TRIM(english)), LOWER(TRIM(afaan_oromo))
);

INSERT INTO words_fts(words_fts) VALUES('rebuild');