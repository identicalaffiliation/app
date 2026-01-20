package tests

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	// "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/identicalaffiliation/app/internal/logger"
	"github.com/identicalaffiliation/app/internal/repository/entity"
	"github.com/identicalaffiliation/app/internal/repository/psql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Init(db *sql.DB) (*psql.Builder, psql.UserRepository) {
	sqlxDB := sqlx.NewDb(db, "postgres")
	qb := psql.NewQueryBuilder()
	postgres := psql.NewPostgres()
	postgres.DB = sqlxDB
	repo := psql.NewUserRepository(postgres, qb, logger.NewLogger())

	return qb, repo
}

func TestCreate(t *testing.T) {
	dbs, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbs.Close()
	id := uuid.New()
	name := "vlad"
	email := "123@mail.ru"
	password := "123"
	user := &entity.User{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}

	qb, repo := Init(dbs)
	query, _, err := qb.Builder.Insert("users").Columns("id", "name", "email", "password").Values(user.ID, user.Name, user.Email, user.Password).Suffix("RETURNING id, created_at").ToSql()

	rows := sqlmock.NewRows([]string{"id", "created_at"}).AddRow(id, time.Now())
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id, name, email, password).
		WillReturnRows(rows)

	err = repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NotNil(t, user.CreatedAt)
	assert.NotNil(t, user.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// func TestGetByID(t *testing.T) {
// 	dbs, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer dbs.Close()

// 	qb, repo := Init(dbs)
// 	time := time.Now()

// 	ids := make(uuid.UUIDs, 2)
// 	for i := range ids {
// 		ids[i] = uuid.New()
// 	}

// 	testTable := []struct {
// 		testName string
// 		ID       uuid.UUID
// 		expected *entity.User
// 	}{
// 		{
// 			testName: "good 1",
// 			ID:       ids[0],
// 			expected: &entity.User{
// 				ID:        ids[0],
// 				Name:      "vlad",
// 				Email:     "123@mail.ru",
// 				Password:  "123123",
// 				CreatedAt: time,
// 				UpdatedAt: time,
// 			},
// 		},
// 	}
// }
// 	for _, testCase := range testTable {
// 		t.Run(testCase.testName, func(t *testing.T) {
// 			sql, _, err := qb.Builder.Select("id, name, email, password, created_at, updated_at").From("users").Where(squirrel.Eq{"id": testCase.ID}).ToSql()
// 			assert.NoError(t, err)

// 			mock.ExpectQuery(regexp.QuoteMeta(sql)).WithArgs(testCase.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"})))

// 			user, err := repo.GetByID(context.Background(), testCase.ID)
// 			assert.NoError(t, err)
// 			assert.Equal(t, testCase.expected.ID, user.ID)
// 		}
// }
// }
