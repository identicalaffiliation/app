package tests

import (
	"database/sql"

	"github.com/identicalaffiliation/app/internal/logger"
	"github.com/identicalaffiliation/app/internal/repository/psql"
	"github.com/jmoiron/sqlx"
)

func InitUser(db *sql.DB) psql.UserRepository {
	sqlxDB := sqlx.NewDb(db, "postgres")
	postgres := psql.NewPostgres()
	postgres.DB = sqlxDB
	repo := psql.NewUserRepository(postgres, logger.NewLogger())

	return repo
}

func InitTodo(db *sql.DB) psql.TodoRepository {
	sqlxDB := sqlx.NewDb(db, "postgres")
	postgres := psql.NewPostgres()
	postgres.DB = sqlxDB
	repo := psql.NewTodoRepository(postgres, logger.NewLogger())

	return repo
}
