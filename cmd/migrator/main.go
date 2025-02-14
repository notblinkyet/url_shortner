package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/notblinkyet/url_shortner/internal/config"
	my_errors "github.com/notblinkyet/url_shortner/internal/errors"
)

func main() {
	var d bool

	cfg := config.MustLoad()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	flag.BoolVar(&d, "d", false, "migrate down")
	flag.Parse()

	pass := os.Getenv("POSTGRES_PASS")
	if pass == "" {
		panic(my_errors.ErrFindPass)
	}
	dbURL := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.MainStorage.Type, cfg.MainStorage.Username, pass, cfg.MainStorage.Host,
		cfg.MainStorage.Port, cfg.MainStorage.Database)
	if cfg.MigrationPath == "" {
		panic("migration_path is required")
	}
	m, err := migrate.New(fmt.Sprintf("file://%s", cfg.MigrationPath), dbURL)
	if err != nil {
		panic(err)
	}

	if d {
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				logger.Println("no migtion to apply down")
				return
			}
			panic(err)
		}
		logger.Println("migrations successfully applied down")
		return
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Println("no migtion to apply up")
			return
		}
		panic(err)
	}
	logger.Println("migrations successfully applied up")
}
