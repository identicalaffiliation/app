package tests

import (
	"database/sql"

	"github.com/identicalaffiliation/app/internal/logger"
	"github.com/identicalaffiliation/app/internal/repository/psql"
	"github.com/jmoiron/sqlx"
)

func InitUser(db *sql.DB) psql.UserRepository {
	sqlxDB := sqlx.NewDb(db, "postgres")
	qb := psql.NewQueryBuilder()
	postgres := psql.NewPostgres()
	postgres.DB = sqlxDB
	repo := psql.NewUserRepository(postgres, qb, logger.NewLogger())

	return repo
}

func InitTodo(db *sql.DB) psql.TodoRepository {
	sqlxDB := sqlx.NewDb(db, "postgres")
	qb := psql.NewQueryBuilder()
	postgres := psql.NewPostgres()
	postgres.DB = sqlxDB
	repo := psql.NewTodoRepository(postgres, qb, logger.NewLogger())

	return repo
}
