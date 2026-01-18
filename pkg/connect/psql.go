package connect

import (
	"fmt"

	"github.com/identicalaffiliation/app/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const POSTGRES_DRIVER string = "postgres"

func toString(cfg *config.AppConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password,
		cfg.Database.Name, cfg.Database.SSLmode)
}

func ConnectToDB(cfg *config.AppConfig) (*sqlx.DB, error) {
	psql, err := sqlx.Open(POSTGRES_DRIVER, toString(cfg))
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	if err := psql.Ping(); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return psql, err
}
