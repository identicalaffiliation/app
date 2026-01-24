package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/identicalaffiliation/app/internal/logger"
	"github.com/identicalaffiliation/app/internal/repository/entity"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *entity.Todo) error
	GetTodosByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Todo, error)
	GetTodoByUserID(ctx context.Context, todoID, userID uuid.UUID) (*entity.Todo, error)
	UpdateStatus(ctx context.Context, newStatus TodoStatus, todoID, userID uuid.UUID) error
	UpdateContent(ctx context.Context, newContent string, todoID, userID uuid.UUID) error
	Delete(ctx context.Context, todoID, userID uuid.UUID) error
}

type todoRepository struct {
	db     *Postgres
	qb     *builder
	logger *logger.Logger
}

func NewTodoRepository(db *Postgres, logger *logger.Logger) TodoRepository {
	qb := NewQueryBuilder()

	return &todoRepository{
		db:     db,
		qb:     qb,
		logger: logger,
	}
}

func (tr *todoRepository) Create(ctx context.Context, todo *entity.Todo) error {
	sql, args, err := tr.qb.Builder.Insert("todos").Columns("id", "user_id", "content",
		"status").Values(todo.ID, todo.UserID, todo.Content, todo.Status).
		Suffix("RETURNING id, created_at").ToSql()
	if err != nil {
		tr.logger.Logger.Error("failed to build query for create todo",
			"operation", "create todo",
			"todo_id", todo.ID.String(),
			"error", err.Error(),
		)

		return ErrFailBuildQuery
	}

	err = tr.db.DB.QueryRowxContext(ctx, sql, args...).Scan(&todo.ID, &todo.CreatedAt)
	if err != nil {
		tr.logger.Logger.Error("failed to create todo",
			"operation", "create todo",
			"todo_id", todo.ID.String(),
			"error", err.Error(),
		)

		return fmt.Errorf("insert todo: %w", err)
	}

	return nil
}

func (tr *todoRepository) GetTodosByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Todo, error) {
	sql, args, err := tr.qb.Builder.Select("id, user_id, content, status, created_at, updated_at").
		From("todos").Where(squirrel.Eq{"user_id": userID}).ToSql()
	if err != nil {
		tr.logger.Logger.Error("failed to build query for get todos",
			"operation", "get todos",
			"user_id", userID.String(),
			"error", err.Error(),
		)

		return nil, ErrFailBuildQuery
	}

	users := make([]*entity.Todo, 0)
	if err := tr.db.DB.SelectContext(ctx, &users, sql, args...); err != nil {
		tr.logger.Logger.Error("failed to get todos",
			"operation", "get todos",
			"user_id", userID.String(),
			"error", err.Error(),
		)

		return nil, fmt.Errorf("select todos: %w", err)
	}

	return users, nil
}

func (tr *todoRepository) GetTodoByUserID(ctx context.Context, todoID, userID uuid.UUID) (*entity.Todo, error) {
	sql, args, err := tr.qb.Builder.Select("id, user_id, content, status, created_at, updated_at").
		From("todos").Where(squirrel.Eq{"user_id": userID}).Where(squirrel.Eq{"id": todoID}).ToSql()
	if err != nil {
		tr.logger.Logger.Error("failed to build query for get todo",
			"operation", "get todo",
			"user_id", userID.String(),
			"todo_id", todoID.String(),
			"error", err.Error(),
		)

		return nil, ErrFailBuildQuery
	}

	var todo entity.Todo
	if err := tr.db.DB.GetContext(ctx, &todo, sql, args...); err != nil {
		tr.logger.Logger.Error("failed to get todo",
			"operation", "get todo",
			"user_id", userID.String(),
			"todo_id", todoID.String(),
			"error", err.Error(),
		)

		return nil, fmt.Errorf("select todo: %w", err)
	}

	return &todo, nil
}

func (tr *todoRepository) UpdateStatus(ctx context.Context, newStatus TodoStatus, todoID, userID uuid.UUID) error {
	sql, args, err := tr.qb.Builder.Update("todos").Set("status", newStatus).Where(squirrel.Eq{"id": todoID}).
		Where(squirrel.Eq{"user_id": userID}).ToSql()
	if err != nil {
		tr.logger.Logger.Error("failed to build query for update todo status",
			"operation", "update status",
			"user_id", userID.String(),
			"todo_id", todoID.String(),
			"req_status", newStatus,
			"error", err.Error(),
		)

		return ErrFailBuildQuery
	}

	result, err := tr.db.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		tr.logger.Logger.Error("failed to update status",
			"operation", "update status",
			"user_id", userID.String(),
			"todo_id", todoID.String(),
			"req_status", newStatus,
			"error", err.Error(),
		)

		return fmt.Errorf("update status: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		tr.logger.Logger.Error("failed to get affected from update todo status",
			"operation", "update status",
			"user_id", userID.String(),
			"todo_id", todoID.String(),
			"req_status", newStatus,
			"error", err.Error(),
		)

		return ErrGetAffected
	}

	if affected == 0 {
		tr.logger.Logger.Error("failed to update status",
			"operation", "update status",
			"user_id", userID.String(),
			"todo_id", todoID.String(),
			"req_status", newStatus,
			"error", errors.New("todo not found").Error(),
		)

		return errors.New("todo not found")
	}

	return nil
}

func (tr *todoRepository) UpdateContent(ctx context.Context, newContent string, todoID, userID uuid.UUID) error {
	sql, args, err := tr.qb.Builder.Update("todos").Set("content", newContent).Where(squirrel.Eq{"id": todoID}).
		Where(squirrel.Eq{"user_id": userID}).ToSql()
	if err != nil {
		tr.logger.Logger.Error("failed to build query for update todo content",
			"operation", "update content",
			"user_id", userID.String(),
			"todo_id", todoID.String(),
			"req_content", newContent,
			"error", err.Error(),
		)

		return ErrFailBuildQuery
	}

	result, err := tr.db.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		tr.logger.Logger.Error("failed to update content",
			"operation", "update content",
			"user_id", userID.String(),
			"todo_id", todoID.String(),
			"req_content", newContent,
			"error", err.Error(),
		)

		return fmt.Errorf("update content: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		tr.logger.Logger.Error("failed to get affected from update todo content",
			"operation", "update content",
			"user_id", userID.String(),
			"todo_id", todoID.String(),
			"req_content", newContent,
			"error", err.Error(),
		)

		return ErrGetAffected
	}

	if affected == 0 {
		tr.logger.Logger.Error("failed to update todo status",
			"operation", "update content",
			"user_id", userID.String(),
			"todo_id", todoID.String(),
			"req_content", newContent,
			"error", errors.New("todo not found").Error(),
		)

		return errors.New("todo not found")
	}

	return nil
}

func (tr *todoRepository) Delete(ctx context.Context, todoID, userID uuid.UUID) error {
	sql, args, err := tr.qb.Builder.Delete("todos").Where(squirrel.Eq{"id": todoID}).
		Where(squirrel.Eq{"user_id": userID}).ToSql()
	if err != nil {
		tr.logger.Logger.Error("failed to build query for delete todo",
			"operation", "delete todo",
			"user_id", userID.String(),
			"todo_id", todoID.String(),
			"error", err.Error(),
		)

		return ErrFailBuildQuery
	}

	result, err := tr.db.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		tr.logger.Logger.Error("failed to delete todo",
			"operation", "delete todo",
			"user_id", userID.String(),
			"todo_id", todoID.String(),
			"error", err.Error(),
		)

		return fmt.Errorf("delete todo: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		tr.logger.Logger.Error("failed to get affected from delete todo",
			"operation", "delete todo",
			"user_id", userID.String(),
			"todo_id", todoID.String(),
			"error", err.Error(),
		)

		return ErrGetAffected
	}

	if affected == 0 {
		tr.logger.Logger.Error("failed to delete todo",
			"operation", "delete todo",
			"user_id", userID.String(),
			"todo_id", todoID.String(),
			"error", errors.New("todo not found").Error(),
		)

		return errors.New("todo not found")
	}

	return nil
}
