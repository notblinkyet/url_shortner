package urlshortner

import (
	"github.com/notblinkyet/url_shortner/internal/lib/generate"
	"github.com/notblinkyet/url_shortner/internal/my_errors"
	"github.com/notblinkyet/url_shortner/internal/repository"
	"github.com/notblinkyet/url_shortner/internal/services"
)

type mockServices struct {
	repo repository.IRepository
}

func NewMockServices() *mockServices {
	return &mockServices{
		repo: services.NewMockRepo(),
	}
}

func (s *mockServices) Create(url string) (string, error) {
	if url == "" {
		return "", my_errors.ErrEmptyURL
	}
	shortUrl, err := generate.Generate(10)
	if err != nil {
		return "", err
	}
	err = s.repo.Create(url, shortUrl)
	if err != nil {
		return "", err
	}
	return shortUrl, err
}

func (s *mockServices) Get(shortUrl string) (string, error) {
	if len(shortUrl) != 10 {
		return "", my_errors.ErrLenShortUrl
	}
	return s.repo.Get(shortUrl)
}
