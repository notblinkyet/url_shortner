package services

import (
	"log"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/notblinkyet/url_shortner/internal/my_errors"
)

func TestCreateGet(t *testing.T) {
	log := log.Default()
	m := NewMockRepo()
	svc := NewServices(log, m)

	example := gofakeit.URL()
	shortUrl, err := svc.Create(example)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	url, err := svc.Get(shortUrl)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if url != example {
		t.Fatalf("expected url: %s\n, got %s", example, url)
	}
}

func TestCreateRepieted(t *testing.T) {
	log := log.Default()
	m := NewMockRepo()
	svc := NewServices(log, m)

	example := gofakeit.URL()
	_, err := svc.Create(example)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	_, err = svc.Create(example)
	if err != my_errors.ErrAlreadyExist {
		t.Fatalf("expected error %v\ngot %v", my_errors.ErrAlreadyExist, err)
	}
}

func TestGetEmpty(t *testing.T) {
	log := log.Default()
	m := NewMockRepo()
	svc := NewServices(log, m)

	example := gofakeit.URL()[:10]
	_, err := svc.Get(example)
	if err != my_errors.ErrAliaceDontUse {
		t.Fatalf("expected error %v\ngot %v", my_errors.ErrAliaceDontUse, err)
	}
}
