package mainstorage

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/notblinkyet/url_shortner/internal/config"
	my_errors "github.com/notblinkyet/url_shortner/internal/my_errors"
)

var (
	ErrUrlAlreadyCreated error = errors.New("url already exists")
	ErrToSave            error = errors.New("can't save")
)

type IMainStorage interface {
	Create(url, shortUrl string) error
	Get(shortUrl string) (string, error)
}

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(cfg config.MainStorage) (*PostgresStorage, error) {
	ctx := context.Background()
	pass := os.Getenv("POSTGRES_PASS")
	if pass == "" {
		return nil, my_errors.ErrFindPass
	}
	pool, err := pgxpool.Connect(ctx,
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.Username, pass, cfg.Host, cfg.Port, cfg.Database))
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{
		pool: pool,
	}, nil
}

func (storage *PostgresStorage) Create(url string, shortUrl string) error {
	ctx := context.Background()
	var a string
	sql := `
		SELECT short_url
		FROM urls
		WHERE short_url=$1
	`
	tx, err := storage.pool.Begin(ctx)
	if err != nil {
		return err
	}
	err = tx.QueryRow(ctx, sql, shortUrl).Scan(&a)
	if err == nil {
		return my_errors.ErrAliasAlreadyUse
	}
	sql = `
		INSERT INTO urls(url, short_url)
		VALUES ($1, $2)
	`
	res, err := storage.pool.Exec(ctx, sql, url, shortUrl)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return my_errors.ErrAlreadyExist
		}
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrToSave
	}
	return nil
}

func (storage *PostgresStorage) Get(shortUrl string) (string, error) {
	ctx := context.Background()
	var url string
	sql := `
		SELECT url
		FROM urls
		WHERE short_url=$1
	`
	err := storage.pool.QueryRow(ctx, sql, shortUrl).Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", my_errors.ErrAliaceDontUse
		}
		return "", err
	}
	return url, nil
}
