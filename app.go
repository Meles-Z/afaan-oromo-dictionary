package main

import (
	"context"
	"log"

	"afaan-oromo-dictionary/internal/db"
	"afaan-oromo-dictionary/internal/dictionary"
)

type App struct {
	ctx     context.Context
	service *dictionary.Service
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	conn, err := db.New("dictionary.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err := db.RunMigrations(conn); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	if err := db.SeedIfEmpty(conn); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	repo := dictionary.NewRepository(conn)
	a.service = dictionary.NewService(repo)
}
func (a *App) SearchWords(query string, direction string) []dictionary.Word {
	words, err := a.service.Search(query, direction)
	if err != nil {
		log.Printf("search error: %v", err)
		return []dictionary.Word{}
	}
	return words
}

func (a *App) GetWord(id int) (*dictionary.Word, error) {
	return a.service.GetWord(id)
}

func (a *App) CreateWord(input dictionary.WordInput) (int, error) {
	return a.service.CreateWord(input)
}

func (a *App) UpdateWord(id int, input dictionary.WordInput) error {
	return a.service.UpdateWord(id, input)
}

func (a *App) DeleteWord(id int) error {
	return a.service.DeleteWord(id)
}

func (a *App) ListCategories() []dictionary.Category {
	cats, err := a.service.ListCategories()
	if err != nil {
		log.Printf("list categories error: %v", err)
		return []dictionary.Category{}
	}
	return cats
}
