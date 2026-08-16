package dictionary

import "strings"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Search(query, direction string) ([]Word, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []Word{}, nil
	}
	return s.repo.Search(strings.ToLower(query), direction)
}

func (s *Service) GetWord(id int) (*Word, error) {
	return s.repo.GetByID(id)
}

func (s *Service) CreateWord(input WordInput) (int, error) {
	if err := validate(input); err != nil {
		return 0, err
	}
	input.English = strings.TrimSpace(input.English)
	input.AfaanOromo = strings.TrimSpace(input.AfaanOromo)
	return s.repo.Create(input)
}

func (s *Service) UpdateWord(id int, input WordInput) error {
	if err := validate(input); err != nil {
		return err
	}
	return s.repo.Update(id, input)
}

func (s *Service) DeleteWord(id int) error {
	return s.repo.Delete(id)
}

func (s *Service) ListCategories() ([]Category, error) {
	return s.repo.ListCategories()
}

func validate(w WordInput) error {
	if strings.TrimSpace(w.English) == "" {
		return ErrEmptyEnglish
	}
	if strings.TrimSpace(w.AfaanOromo) == "" {
		return ErrEmptyAfaanOromo
	}
	return nil
}