package repository

import (
	my_errors "github.com/notblinkyet/url_shortner/internal/my_errors"

	"github.com/notblinkyet/url_shortner/internal/repository/cache"
	mainstorage "github.com/notblinkyet/url_shortner/internal/repository/main_storage"
)

type IRepository interface {
	Create(url, shortUrl string) error
	Get(shortUrl string) (string, error)
}

type Repository struct {
	cache       cache.ICache
	mainStorage mainstorage.IMainStorage
}

func NewRepository(cache cache.ICache, mainStorage mainstorage.IMainStorage) *Repository {
	return &Repository{
		cache:       cache,
		mainStorage: mainStorage,
	}
}

func (repo *Repository) Create(url, shortUrl string) error {
	_, err := repo.cache.Get(url)
	if err == nil {
		return my_errors.ErrAlreadyExist
	}
	err = repo.mainStorage.Create(url, shortUrl)
	if err != nil {
		return err
	}
	return repo.cache.Set(url, shortUrl)
}

func (repo *Repository) Get(shortUrl string) (string, error) {
	url, err := repo.cache.Get(shortUrl)
	if err == nil {
		return url, err
	}
	return repo.mainStorage.Get(shortUrl)
}
