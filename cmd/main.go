package main

import (
	"github.com/identicalaffiliation/app/internal/config"
	"github.com/identicalaffiliation/app/internal/logger"
	"github.com/identicalaffiliation/app/internal/repository/psql"
	"github.com/identicalaffiliation/app/internal/service"
	"github.com/identicalaffiliation/app/pkg/parse"
)

func main() {
	path := parse.FlagInit()

	cfg := config.MustLoadConfig(path)

	logger := logger.NewLogger()

	db := psql.NewPostgres()
	defer db.Close()

	db.MustInit(cfg)

	userRepo := psql.NewUserRepository(db, logger)
	_ = service.NewUserService(userRepo)
	_ = service.NewAuthService(userRepo, cfg.JWTSecret)
}
