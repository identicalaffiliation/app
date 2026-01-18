package main

import (
	"github.com/identicalaffiliation/app/internal/config"
	"github.com/identicalaffiliation/app/internal/repository/psql"
	"github.com/identicalaffiliation/app/pkg/parse"
)

func main() {
	path := parse.FlagInit()

	cfg := config.MustLoadConfig(path)

	db := psql.NewPostgres()
	db.MustInitDB(cfg)
}
