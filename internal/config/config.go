package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Grpc struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type MainStorage struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"db"`
	Username string `yaml:"username"`
}

type Cache struct {
	Type string        `yaml:"type"`
	Host string        `yaml:"host"`
	Port int           `yaml:"port"`
	DB   int           `yaml:"db"`
	Exp  time.Duration `yaml:"exp"`
}

type Config struct {
	Env           string      `yaml:"env"`
	MigrationPath string      `yaml:"migration_path"`
	Grpc          Grpc        `yaml:"grpc"`
	MainStorage   MainStorage `yaml:"main_storage"`
	Cache         Cache       `yaml:"cache"`
}

func MustLoad() *Config {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		panic("can't find CONFIG_PATH in env")
	}
	cfg, err := Load(path)
	if err != nil {
		panic(err)
	}
	return cfg
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var cfg Config
	err = yaml.NewDecoder(file).Decode(&cfg)
	return &cfg, err
}
