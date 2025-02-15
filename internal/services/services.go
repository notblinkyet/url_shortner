package services

import (
	"errors"
	"log"

	"github.com/notblinkyet/url_shortner/internal/lib/generate"
	my_errors "github.com/notblinkyet/url_shortner/internal/my_errors"
	"github.com/notblinkyet/url_shortner/internal/repository"
)

var ()

type IServices interface {
	Create(url string) (string, error)
	Get(shortUrl string) (string, error)
}

type Services struct {
	logger *log.Logger
	repo   repository.IRepository
}

func NewServices(logger *log.Logger, repo repository.IRepository) *Services {
	return &Services{
		logger: logger,
		repo:   repo,
	}
}

func (services *Services) Create(url string) (string, error) {
	if len(url) == 0 {
		services.logger.Println(my_errors.ErrEmptyURL)
		return "", my_errors.ErrEmptyURL
	}
	err := my_errors.ErrAliasAlreadyUse
	var shortUrl string
	for errors.Is(err, my_errors.ErrAliasAlreadyUse) {
		shortUrl, err = generate.Generate(10)
		if err != nil {
			return "", err
		}
		err = services.repo.Create(url, string(shortUrl))
		if err == nil {
			break
		}
	}
	if err != nil {
		services.logger.Println(err)
	} else {
		services.logger.Println("successfully create shortUrl")
	}
	return shortUrl, err
}

func (services *Services) Get(shortUrl string) (string, error) {
	if len(shortUrl) != 10 {
		services.logger.Println(my_errors.ErrLenShortUrl)
		return "", my_errors.ErrLenShortUrl
	}
	url, err := services.repo.Get(shortUrl)
	if err != nil {
		services.logger.Println(err)
	} else {
		services.logger.Println("successfully get url from short")
	}
	return url, err
}
