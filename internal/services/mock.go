package services

import my_errors "github.com/notblinkyet/url_shortner/internal/my_errors"

type mockRepository struct {
	data map[string]string
}

func NewMockRepo() *mockRepository {
	return &mockRepository{
		data: make(map[string]string),
	}
}

func (m *mockRepository) Create(url, shortUrl string) error {
	if _, ok := m.data[url]; ok {
		return my_errors.ErrAlreadyExist
	}
	m.data[url] = shortUrl
	return nil
}

func (m *mockRepository) Get(shortUrl string) (string, error) {
	for key, val := range m.data {
		if shortUrl == val {
			return key, nil
		}
	}
	return "", my_errors.ErrAliaceDontUse
}
