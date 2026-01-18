package psql

import (
	"github.com/Masterminds/squirrel"
	"github.com/identicalaffiliation/app/internal/config"
	"github.com/identicalaffiliation/app/pkg/connect"
	"github.com/jmoiron/sqlx"
)

type Builder struct {
	builder squirrel.StatementBuilderType
}

func newQueryBuilder() *Builder {
	return &Builder{
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

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
