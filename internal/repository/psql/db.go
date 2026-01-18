package psql

import (
	"github.com/identicalaffiliation/app/internal/config"
	"github.com/identicalaffiliation/app/pkg/connect"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres() *Postgres {
	return &Postgres{}
}

func (p *Postgres) MustInitDB(cfg *config.AppConfig) {
	db, err := connect.ConnectToDB(cfg)
	if err != nil {
		panic(err)
	}

	p.db = db
}
