package services

import (
	"crypto/rand"
	"errors"
	"log"
	"math/big"

	my_errors "github.com/notblinkyet/url_shortner/internal/errors"
	"github.com/notblinkyet/url_shortner/internal/repository"
)

var (
	ErrEmptyURL    error = errors.New("empty url")
	ErrLenShortUrl error = errors.New("length or short url != 10")
)

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
		services.logger.Println(ErrEmptyURL)
	}
	var shortUrl []byte = make([]byte, 10)
	var set string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"
	err := my_errors.ErrAliasAlreadyUse
	for errors.Is(err, my_errors.ErrAliasAlreadyUse) {
		for i := 0; i < 10; i++ {
			r, err := rand.Int(rand.Reader, big.NewInt(63))
			if err != nil {
				return "", err
			}
			shortUrl[i] = set[r.Int64()]
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
	return string(shortUrl), err
}

func (services *Services) Get(shortUrl string) (string, error) {
	if len(shortUrl) != 10 {
		services.logger.Println(ErrLenShortUrl)
	}
	url, err := services.repo.Get(shortUrl)
	if err != nil {
		services.logger.Println(err)
	} else {
		services.logger.Println("successfully get url from short")
	}
	return url, err
}
