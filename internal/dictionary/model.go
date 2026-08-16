package dictionary

import "time"

type Word struct {
	ID            int       `json:"id"`
	English       string    `json:"english"`
	AfaanOromo    string    `json:"afaanOromo"`
	PartOfSpeech  string    `json:"partOfSpeech"`
	ExampleEn     string    `json:"exampleEn"`
	ExampleOm     string    `json:"exampleOm"`
	Pronunciation string    `json:"pronunciation"`
	Categories    []Category `json:"categories,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// WordInput is what the frontend sends when creating/updating a word.
// Kept separate from Word so the frontend never has to worry about
// ID, timestamps, or nested categories on write.
type WordInput struct {
	English       string `json:"english"`
	AfaanOromo    string `json:"afaanOromo"`
	PartOfSpeech  string `json:"partOfSpeech"`
	ExampleEn     string `json:"exampleEn"`
	ExampleOm     string `json:"exampleOm"`
	Pronunciation string `json:"pronunciation"`
	CategoryIDs   []int  `json:"categoryIds"`
}