package app

import (
	"log"

	grpcapp "github.com/notblinkyet/url_shortner/internal/app/grpc"
	"github.com/notblinkyet/url_shortner/internal/config"
	"github.com/notblinkyet/url_shortner/internal/repository"
	"github.com/notblinkyet/url_shortner/internal/repository/cache"
	mainstorage "github.com/notblinkyet/url_shortner/internal/repository/main_storage"
	"github.com/notblinkyet/url_shortner/internal/services"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *log.Logger, cfg *config.Config) *App {

	var (
		mainStorage  mainstorage.IMainStorage
		cacheStorage cache.ICache
		err          error
	)
	switch cfg.MainStorage.Type {
	case "postgres":
		mainStorage, err = mainstorage.NewPostgresStorage(cfg.MainStorage)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Unknown DB type")
	}

	switch cfg.Cache.Type {
	case "redis":
		cacheStorage, err = cache.NewRedisCache(cfg.Cache)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Unknown cache type")
	}
	layerRepo := repository.NewRepository(cacheStorage, mainStorage)
	layerServer := services.NewServices(log, layerRepo)
	GRPCServer := grpcapp.New(log, layerServer, cfg.Grpc.Port, cfg.Grpc.Host, cfg.Grpc.Timeout)
	return &App{
		GRPCServer: GRPCServer,
	}
}
