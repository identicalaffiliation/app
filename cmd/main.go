package main

import (
	"github.com/identicalaffiliation/app/internal/config"
	"github.com/identicalaffiliation/app/internal/logger"
	"github.com/identicalaffiliation/app/internal/repository/psql"
	"github.com/identicalaffiliation/app/pkg/parse"
)

func main() {
	path := parse.FlagInit()

	cfg := config.MustLoadConfig(path)

	logger := logger.NewLogger()

	db := psql.NewPostgres()
	defer db.Close()

	queryBuilder := psql.NewQueryBuilder()

	db.MustInitDB(cfg)

	_ = psql.NewUserRepository(db, queryBuilder, logger)
}
