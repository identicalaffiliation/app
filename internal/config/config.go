package config

import (
	"errors"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var (
	ErrInvalidENV    error = errors.New("env must be load in app")
	ErrInvalidConfig error = errors.New("invalid config file")
)

type PostgresConfig struct {
	Name     string `env:"DB_NAME"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	SSLmode  string `env:"DB_SSLMODE"`
}

type HTTPConfig struct {
	Port string `yaml:"http_port"`
}

type AppConfig struct {
	Database   PostgresConfig
	HTTPServer HTTPConfig `yaml:"http"`
}

func MustLoadConfig(path string) *AppConfig {
	cfg := &AppConfig{}

	godotenv.Load()

	if err := cleanenv.ReadEnv(cfg); err != nil {
		panic(ErrInvalidENV)
	}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic(ErrInvalidConfig)
	}

	return cfg
}
