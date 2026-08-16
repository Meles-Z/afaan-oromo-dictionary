package dictionary

import "errors"

var (
	ErrNotFound      = errors.New("word not found")
	ErrEmptyEnglish  = errors.New("english field is required")
	ErrEmptyAfaanOromo = errors.New("afaan oromo field is required")
)