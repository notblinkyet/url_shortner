package cache

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/notblinkyet/url_shortner/internal/config"
	"github.com/redis/go-redis/v9"
)

var (
	ErrFindPass  error = errors.New("can't find REDIS_PASS in env")
	ErrToconnect error = errors.New("can't connect to Redis")
)

type ICache interface {
	Get(url string) (string, error)
	Set(url, shortUrl string) error
}

type RedisCache struct {
	client *redis.Client
	exp    time.Duration
}

func NewRedisCache(cfg config.Cache) (*RedisCache, error) {
	ctx := context.Background()
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	pass := os.Getenv("REDIS_PASS")
	if pass == "" {
		return nil, ErrFindPass
	}
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       cfg.DB,
		Password: pass,
	})
	if client == nil {
		return nil, ErrToconnect
	}

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &RedisCache{
		client: client,
		exp:    cfg.Exp,
	}, nil
}

func (cache *RedisCache) Get(url string) (string, error) {
	ctx := context.Background()
	res := cache.client.Get(ctx, url)
	if err := res.Err(); err != nil {
		return "", err
	}
	return res.Val(), nil
}

func (cache *RedisCache) Set(url, shortUrl string) error {
	ctx := context.Background()
	if err := cache.client.Set(ctx, url, shortUrl, cache.exp).Err(); err != nil {
		return err
	}
	return cache.client.Set(ctx, shortUrl, url, cache.exp).Err()
}
