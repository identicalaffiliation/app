package psql

import (
	"github.com/identicalaffiliation/app/internal/config"
	"github.com/identicalaffiliation/app/pkg/connect"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	DB *sqlx.DB
}

func NewPostgres() *Postgres {
	return &Postgres{}
}

func (p *Postgres) Close() error {
	return p.DB.Close()
}

func (p *Postgres) MustInitUserDB(cfg *config.AppConfig) {
	db, err := connect.ConnectToDB(cfg)
	if err != nil {
		panic(err)
	}

	p.DB = db
}
