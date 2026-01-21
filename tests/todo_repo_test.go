package tests

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/identicalaffiliation/app/internal/repository/entity"
	"github.com/identicalaffiliation/app/internal/repository/psql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTodo(t *testing.T) {
	type testCase struct {
		testName  string
		mockSetup func(mock sqlmock.Sqlmock, id, user_id uuid.UUID, content string, status psql.TodoStatus)
		inputTodo *entity.Todo
		expected  *entity.Todo
	}

	testTime := time.Now()
	todoID := uuid.New()
	userID := uuid.New()

	testTable := []testCase{
		{
			testName: "success – todo created",
			mockSetup: func(mock sqlmock.Sqlmock, id, user_id uuid.UUID, content string, status psql.TodoStatus) {
				rows := sqlmock.NewRows([]string{"id", "created_at"}).AddRow(todoID, testTime)

				mock.ExpectQuery(regexp.QuoteMeta(TODO_CREATE_QUERY)).WithArgs(id, user_id, content, status).WillReturnRows(rows)
			},
			inputTodo: &entity.Todo{
				ID:      todoID,
				UserID:  userID,
				Content: "breakfast",
				Status:  "todo",
			},
			expected: &entity.Todo{
				ID:        todoID,
				UserID:    userID,
				Content:   "breakfast",
				Status:    "todo",
				CreatedAt: testTime,
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			repo := InitTodo(db)

			testCase.mockSetup(mock, testCase.expected.ID, testCase.expected.UserID,
				testCase.expected.Content, psql.TodoStatus(testCase.expected.Status))

			err = repo.Create(context.Background(), testCase.inputTodo)
			require.NoError(t, err)
			assert.Equal(t, testCase.expected.ID, testCase.inputTodo.ID)
			assert.Equal(t, testCase.expected.CreatedAt, testCase.inputTodo.CreatedAt)

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetTodos(t *testing.T) {
	type testCase struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock, userID uuid.UUID)
		userID        uuid.UUID
		expectedTodos []*entity.Todo
	}

	testTime := time.Now()
	testTimea := time.Now()
	todoID := uuid.New()
	todoIDa := uuid.New()
	userID := uuid.New()
	userIDa := uuid.New()
	content := "breakfast"
	contenta := "running"
	status := psql.Todo
	statusa := psql.Done

	testTable := []testCase{
		{
			testName: "success – todos found",
			mockSetup: func(mock sqlmock.Sqlmock, userID uuid.UUID) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "content", "status", "created_at", "updated_at"}).
					AddRow(todoID, userID, content, status, testTime, testTime).
					AddRow(todoIDa, userIDa, contenta, statusa, testTimea, testTimea)

				mock.ExpectQuery(TODO_GET_TODOS_BY_USER_ID).WithArgs(userID).WillReturnRows(rows)
			},
			userID: userID,
			expectedTodos: []*entity.Todo{
				{
					ID:        todoID,
					UserID:    userID,
					Content:   content,
					Status:    string(status),
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
				{
					ID:        todoIDa,
					UserID:    userIDa,
					Content:   contenta,
					Status:    string(statusa),
					CreatedAt: testTimea,
					UpdatedAt: testTimea,
				},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			repo := InitTodo(db)

			testCase.mockSetup(mock, testCase.userID)

			result, err := repo.GetTodosByUserID(context.Background(), testCase.userID)
			require.NoError(t, err)

			assert.NotNil(t, result)
			assert.Equal(t, testCase.expectedTodos, result)

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetTodoByID(t *testing.T) {
	type testCase struct {
		testName     string
		mockSetup    func(mock sqlmock.Sqlmock, userID, todoID uuid.UUID)
		userID       uuid.UUID
		todoID       uuid.UUID
		expectedTodo *entity.Todo
	}

	testTime := time.Now()
	todoID := uuid.New()
	userID := uuid.New()
	content := "breakfast"
	status := psql.Todo

	testTable := []testCase{
		{
			testName: "success – todos found",
			mockSetup: func(mock sqlmock.Sqlmock, userID, todoID uuid.UUID) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "content", "status", "created_at", "updated_at"}).
					AddRow(todoID, userID, content, status, testTime, testTime)

				mock.ExpectQuery(regexp.QuoteMeta(TODO_GET_TODO_BY_USER_ID)).WithArgs(userID, todoID).WillReturnRows(rows)

			},
			userID: userID,
			todoID: todoID,
			expectedTodo: &entity.Todo{
				ID:        todoID,
				UserID:    userID,
				Content:   content,
				Status:    string(status),
				CreatedAt: testTime,
				UpdatedAt: testTime,
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			repo := InitTodo(db)

			testCase.mockSetup(mock, testCase.userID, testCase.todoID)

			result, err := repo.GetTodoByUserID(context.Background(), testCase.todoID, testCase.userID)
			require.NoError(t, err)

			assert.NotNil(t, result)
			assert.Equal(t, testCase.expectedTodo.ID, result.ID)
			assert.Equal(t, testCase.expectedTodo.Content, result.Content)
			assert.Equal(t, testCase.expectedTodo.CreatedAt, result.CreatedAt)

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateStatus(t *testing.T) {
	type testCase struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock, userID, todoID uuid.UUID, status psql.TodoStatus)
		userID        uuid.UUID
		todoID        uuid.UUID
		status        psql.TodoStatus
		expectedError string
	}

	todoID := uuid.New()
	userID := uuid.New()
	status := psql.Todo
	e := errors.New("todo not found").Error()

	testTable := []testCase{
		{
			testName: "success – todos found",
			mockSetup: func(mock sqlmock.Sqlmock, userID, todoID uuid.UUID, status psql.TodoStatus) {

				mock.ExpectExec(regexp.QuoteMeta(TODO_UPDATE_STATUS)).WithArgs(status, todoID, userID).WillReturnResult(sqlmock.NewResult(0, 1))

			},
			userID:        userID,
			todoID:        todoID,
			status:        status,
			expectedError: "",
		},
		{
			testName: "error – todo not found",
			mockSetup: func(mock sqlmock.Sqlmock, userID, todoID uuid.UUID, status psql.TodoStatus) {
				mock.ExpectExec(regexp.QuoteMeta(TODO_UPDATE_STATUS)).WithArgs(status, todoID, userID).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			userID:        userID,
			todoID:        todoID,
			status:        status,
			expectedError: e,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			repo := InitTodo(db)

			if testCase.expectedError != "" {
				testCase.mockSetup(mock, testCase.userID, testCase.todoID, status)
				err := repo.UpdateStatus(context.Background(), testCase.status, testCase.todoID, testCase.userID)
				require.Error(t, err)
				assert.Equal(t, err.Error(), testCase.expectedError)

				require.NoError(t, mock.ExpectationsWereMet())
			} else {
				testCase.mockSetup(mock, testCase.userID, testCase.todoID, status)
				err := repo.UpdateStatus(context.Background(), testCase.status, testCase.todoID, testCase.userID)
				require.NoError(t, err)

				require.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}

func TestUpdateContent(t *testing.T) {
	type testCase struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock, userID, todoID uuid.UUID, content string)
		userID        uuid.UUID
		todoID        uuid.UUID
		content       string
		expectedError string
	}

	todoID := uuid.New()
	userID := uuid.New()
	content := "breakfast"
	e := errors.New("todo not found").Error()

	testTable := []testCase{
		{
			testName: "success – todo updated",
			mockSetup: func(mock sqlmock.Sqlmock, userID, todoID uuid.UUID, content string) {

				mock.ExpectExec(regexp.QuoteMeta(TODO_UPDATE_CONTENT)).WithArgs(content, todoID, userID).WillReturnResult(sqlmock.NewResult(0, 1))

			},
			userID:        userID,
			todoID:        todoID,
			content:       content,
			expectedError: "",
		},
		{
			testName: "error – todo not found",
			mockSetup: func(mock sqlmock.Sqlmock, userID, todoID uuid.UUID, content string) {
				mock.ExpectExec(regexp.QuoteMeta(TODO_UPDATE_CONTENT)).WithArgs(content, todoID, userID).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			userID:        userID,
			todoID:        todoID,
			content:       content,
			expectedError: e,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			repo := InitTodo(db)

			if testCase.expectedError != "" {
				testCase.mockSetup(mock, testCase.userID, testCase.todoID, testCase.content)
				err := repo.UpdateContent(context.Background(), testCase.content, testCase.todoID, testCase.userID)
				require.Error(t, err)
				assert.Equal(t, err.Error(), testCase.expectedError)

				require.NoError(t, mock.ExpectationsWereMet())
			} else {
				testCase.mockSetup(mock, testCase.userID, testCase.todoID, testCase.content)
				err := repo.UpdateContent(context.Background(), testCase.content, testCase.todoID, testCase.userID)
				require.NoError(t, err)

				require.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	type testCase struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock, userID, todoID uuid.UUID)
		userID        uuid.UUID
		todoID        uuid.UUID
		expectedError string
	}

	todoID := uuid.New()
	userID := uuid.New()
	e := errors.New("todo not found").Error()

	testTable := []testCase{
		{
			testName: "success – todo deleted",
			mockSetup: func(mock sqlmock.Sqlmock, userID, todoID uuid.UUID) {

				mock.ExpectExec(regexp.QuoteMeta(TODO_DELETE)).WithArgs(todoID, userID).WillReturnResult(sqlmock.NewResult(0, 1))

			},
			userID:        userID,
			todoID:        todoID,
			expectedError: "",
		},
		{
			testName: "error – todo not found",
			mockSetup: func(mock sqlmock.Sqlmock, userID, todoID uuid.UUID) {
				mock.ExpectExec(regexp.QuoteMeta(TODO_DELETE)).WithArgs(todoID, userID).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			userID:        userID,
			todoID:        todoID,
			expectedError: e,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			repo := InitTodo(db)

			if testCase.expectedError != "" {
				testCase.mockSetup(mock, testCase.userID, testCase.todoID)
				err := repo.Delete(context.Background(), testCase.todoID, testCase.userID)
				require.Error(t, err)
				assert.Equal(t, err.Error(), testCase.expectedError)

				require.NoError(t, mock.ExpectationsWereMet())
			} else {
				testCase.mockSetup(mock, testCase.userID, testCase.todoID)
				err := repo.Delete(context.Background(), testCase.todoID, testCase.userID)
				require.NoError(t, err)

				require.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}
