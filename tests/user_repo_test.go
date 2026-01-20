package tests

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/app/internal/logger"
	"github.com/identicalaffiliation/app/internal/repository/entity"
	"github.com/identicalaffiliation/app/internal/repository/psql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Init(db *sql.DB) psql.UserRepository {
	sqlxDB := sqlx.NewDb(db, "postgres")
	qb := psql.NewQueryBuilder()
	postgres := psql.NewPostgres()
	postgres.DB = sqlxDB
	repo := psql.NewUserRepository(postgres, qb, logger.NewLogger())

	return repo
}

func TestCreateUser(t *testing.T) {
	type testCase struct {
		testName     string
		setupMock    func(mock sqlmock.Sqlmock)
		inputUser    *entity.User
		expectedUser *entity.User
	}

	testTime := time.Now()
	ID1 := uuid.New()

	testCases := []testCase{
		{
			testName: "success – user created",
			setupMock: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO users \(id,name,email,password\) VALUES \(\$1,\$2,\$3,\$4\) RETURNING id, created_at`

				rows := sqlmock.NewRows([]string{"id", "created_at"}).AddRow(ID1, testTime)

				mock.ExpectQuery(query).WithArgs(ID1, "vlad", "123@mail.ru", "123123").WillReturnRows(rows)
			},

			inputUser: &entity.User{
				ID:       ID1,
				Email:    "123@mail.ru",
				Name:     "vlad",
				Password: "123123",
			},
			expectedUser: &entity.User{
				ID:        ID1,
				Email:     "123@mail.ru",
				Name:      "vlad",
				Password:  "123123",
				CreatedAt: testTime,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			repo := Init(db)

			testCase.setupMock(mock)

			err = repo.Create(context.Background(), testCase.inputUser)
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedUser.ID, testCase.inputUser.ID)
			assert.Equal(t, testCase.expectedUser.CreatedAt, testCase.inputUser.CreatedAt)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	type testCase struct {
		testName      string
		setupMock     func(mock sqlmock.Sqlmock, repo *psql.UserRepository)
		inputID       uuid.UUID
		expectedUsers []*entity.User
	}

	testTime := time.Now()
	validID := uuid.New()
	validID2 := uuid.New()

	testCases := []testCase{
		{
			testName: "success – users found",
			setupMock: func(mock sqlmock.Sqlmock, repo *psql.UserRepository) {
				query := `SELECT id, name, email, password, created_at, updated_at FROM users`

				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow(validID, "vlad", "123@mail.ru", "123123", testTime, testTime).
					AddRow(validID2, "ruslan", "321@gmail.com", "321321", testTime, testTime)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
			},
			expectedUsers: []*entity.User{
				{
					ID:        validID,
					Name:      "vlad",
					Email:     "123@mail.ru",
					Password:  "123123",
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
				{
					ID:        validID2,
					Name:      "ruslan",
					Email:     "321@gmail.com",
					Password:  "321321",
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := Init(db)

			testCase.setupMock(mock, &repo)

			users, err := repo.GetAllUsers(context.Background())
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedUsers, users)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetByID(t *testing.T) {
	type testCase struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock, expected *entity.User)
		inputID       uuid.UUID
		expectedUser  *entity.User
		expectedError error
	}

	testTime := time.Now()
	testID := uuid.New()
	testInvalidID := uuid.New()

	testCases := []testCase{
		{
			testName: "success – user found",
			mockSetup: func(mock sqlmock.Sqlmock, expected *entity.User) {
				query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1`

				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow(expected.ID, expected.Name, expected.Email, expected.Password,
						expected.CreatedAt, expected.UpdatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(expected.ID).WillReturnRows(rows)
			},
			inputID: testID,
			expectedUser: &entity.User{
				ID:        testID,
				Name:      "vlad",
				Email:     "123@mail.ru",
				Password:  "123",
				CreatedAt: testTime,
				UpdatedAt: testTime,
			},
			expectedError: nil,
		},
		{
			testName: "error – invalid user ID",
			mockSetup: func(mock sqlmock.Sqlmock, expected *entity.User) {
				query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1`

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(expected.ID).WillReturnError(psql.ErrInvalidUserID)
			},
			inputID:       testInvalidID,
			expectedUser:  nil,
			expectedError: psql.ErrInvalidUserID,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := Init(db)

			if testCase.expectedError != nil {
				testCase.mockSetup(mock, &entity.User{ID: testInvalidID})
				result, err := repo.GetByID(context.Background(), testCase.inputID)
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, err, psql.ErrInvalidUserID)

				require.NoError(t, mock.ExpectationsWereMet())
			} else {
				testCase.mockSetup(mock, testCase.expectedUser)
				result, err := repo.GetByID(context.Background(), testCase.inputID)
				require.NoError(t, err)

				assert.Equal(t, testCase.expectedUser.ID, result.ID)
				assert.Equal(t, testCase.expectedUser.Name, result.Name)
				assert.Equal(t, testCase.expectedUser.CreatedAt, result.CreatedAt)
				require.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}

func TestChangeName(t *testing.T) {
	type testCase struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock, id uuid.UUID)
		inputID       uuid.UUID
		inputName     string
		expectedError error
	}

	validID := uuid.New()
	validName := "a"
	invalidID := uuid.New()
	testTime := time.Now()

	testCases := []testCase{
		{
			testName: "success – name updated",
			mockSetup: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				query := `UPDATE users SET name = \$1 WHERE id = \$2`

				_ = sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow(id, "vlad", "123@mail.ru", "123", testTime, testTime)

				mock.ExpectExec(query).WithArgs("a", id).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			inputID:       validID,
			inputName:     validName,
			expectedError: nil,
		},
		{
			testName: "error – invalid user ID",
			mockSetup: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				query := `UPDATE users SET name = \$1 WHERE id = \$2`

				mock.ExpectExec(query).WithArgs("b", id).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			inputID:       invalidID,
			inputName:     "b",
			expectedError: psql.ErrInvalidUserID,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := Init(db)

			if testCase.expectedError != nil {
				testCase.mockSetup(mock, testCase.inputID)

				err := repo.ChangeName(context.Background(), testCase.inputName, testCase.inputID)
				require.Error(t, err)
				assert.ErrorIs(t, testCase.expectedError, err)
				require.NoError(t, mock.ExpectationsWereMet())
			} else {
				testCase.mockSetup(mock, testCase.inputID)
				err = repo.ChangeName(context.Background(), validName, validID)
				require.NoError(t, err)
				require.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}

func TestChangeEmail(t *testing.T) {
	type testCase struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock, id uuid.UUID)
		inputID       uuid.UUID
		inputEmail    string
		expectedError error
	}

	validID := uuid.New()
	validEmail := "a"
	invalidID := uuid.New()
	testTime := time.Now()

	testCases := []testCase{
		{
			testName: "success – email updated",
			mockSetup: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				query := `UPDATE users SET email = \$1 WHERE id = \$2`

				_ = sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow(id, "vlad", "123@mail.ru", "123", testTime, testTime)

				mock.ExpectExec(query).WithArgs("a", id).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			inputID:       validID,
			inputEmail:    validEmail,
			expectedError: nil,
		},
		{
			testName: "error – invalid user ID",
			mockSetup: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				query := `UPDATE users SET email = \$1 WHERE id = \$2`

				mock.ExpectExec(query).WithArgs("b", id).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			inputID:       invalidID,
			inputEmail:    "b",
			expectedError: psql.ErrInvalidUserID,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := Init(db)

			if testCase.expectedError != nil {
				testCase.mockSetup(mock, testCase.inputID)

				err := repo.ChangeEmail(context.Background(), testCase.inputEmail, testCase.inputID)
				require.Error(t, err)
				assert.ErrorIs(t, testCase.expectedError, err)
				require.NoError(t, mock.ExpectationsWereMet())
			} else {
				testCase.mockSetup(mock, testCase.inputID)
				err = repo.ChangeEmail(context.Background(), validEmail, validID)
				require.NoError(t, err)
				require.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}

func TestChangePassword(t *testing.T) {
	type testCase struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock, id uuid.UUID)
		inputID       uuid.UUID
		inputPassword string
		expectedError error
	}

	validID := uuid.New()
	validPassword := "a"
	invalidID := uuid.New()
	testTime := time.Now()

	testCases := []testCase{
		{
			testName: "success – password updated",
			mockSetup: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				query := `UPDATE users SET password = \$1 WHERE id = \$2`

				_ = sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow(id, "vlad", "123@mail.ru", "123", testTime, testTime)

				mock.ExpectExec(query).WithArgs("a", id).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			inputID:       validID,
			inputPassword: validPassword,
			expectedError: nil,
		},
		{
			testName: "error – invalid user ID",
			mockSetup: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				query := `UPDATE users SET password = \$1 WHERE id = \$2`

				mock.ExpectExec(query).WithArgs("b", id).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			inputID:       invalidID,
			inputPassword: "b",
			expectedError: psql.ErrInvalidUserID,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := Init(db)

			if testCase.expectedError != nil {
				testCase.mockSetup(mock, testCase.inputID)

				err := repo.ChangePassword(context.Background(), testCase.inputPassword, testCase.inputID)
				require.Error(t, err)
				assert.ErrorIs(t, testCase.expectedError, err)
				require.NoError(t, mock.ExpectationsWereMet())
			} else {
				testCase.mockSetup(mock, testCase.inputID)
				err = repo.ChangePassword(context.Background(), validPassword, validID)
				require.NoError(t, err)
				require.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type testCase struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock, id uuid.UUID)
		inputID       uuid.UUID
		ExpectedError error
	}

	testTime := time.Now()
	validID := uuid.New()
	invalidID := uuid.New()

	testCases := []testCase{
		{
			testName: "success – user deleted",
			mockSetup: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				query := `DELETE FROM users WHERE id = \$1`

				_ = sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow(id, "vlad", "123@mail.ru", "123", testTime, testTime)

				mock.ExpectExec(query).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			inputID:       validID,
			ExpectedError: nil,
		},

		{
			testName: "error – invalid id",
			mockSetup: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				query := `DELETE FROM users WHERE id = \$1`

				mock.ExpectExec(query).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			inputID:       invalidID,
			ExpectedError: psql.ErrInvalidUserID,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := Init(db)

			if testCase.ExpectedError != nil {
				testCase.mockSetup(mock, testCase.inputID)
				err := repo.Delete(context.Background(), testCase.inputID)
				require.Error(t, err)
				assert.Equal(t, testCase.ExpectedError, err)
				require.NoError(t, mock.ExpectationsWereMet())
			} else {
				testCase.mockSetup(mock, testCase.inputID)
				err := repo.Delete(context.Background(), testCase.inputID)
				require.NoError(t, err)
				require.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}
