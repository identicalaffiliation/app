package psql

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/identicalaffiliation/app/internal/logger"
	"github.com/identicalaffiliation/app/internal/repository/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
	ChangeName(ctx context.Context, newName string, userID uuid.UUID) error
	ChangeEmail(ctx context.Context, newEmail string, userID uuid.UUID) error
	ChangePassword(ctx context.Context, newPassword string, userID uuid.UUID) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

type userRepository struct {
	db     *Postgres
	qb     *Builder
	logger *logger.Logger
}

func NewUserRepository(db *Postgres, qb *Builder, logger *logger.Logger) UserRepository {
	return &userRepository{
		db:     db,
		qb:     qb,
		logger: logger,
	}
}

func (ur *userRepository) Create(ctx context.Context, user *entity.User) error {
	sql, args, err := ur.qb.Builder.Insert("users").Columns("id", "name", "email",
		"password").Values(user.ID, user.Name, user.Email, user.Password).
		Suffix("RETURNING id, created_at").ToSql()
	if err != nil {
		return ErrFailBuildQuery
	}

	err = ur.db.DB.QueryRowxContext(ctx, sql, args...).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

func (ur *userRepository) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	sql, args, err := ur.qb.Builder.Select("id", "name", "email", "password",
		"created_at", "updated_at").From("users").OrderBy("email").ToSql()
	if err != nil {
		return nil, ErrFailBuildQuery
	}

	var users []*entity.User
	if err := ur.db.DB.SelectContext(ctx, &users, sql, args...); err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}

	return users, nil
}

func (ur *userRepository) GetByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	sql, args, err := ur.qb.Builder.Select("id, name, email, password, created_at, updated_at").
		From("users").Where(squirrel.Eq{"id": userID}).ToSql()
	if err != nil {
		return nil, ErrFailBuildQuery
	}

	var user entity.User
	if err := ur.db.DB.GetContext(ctx, &user, sql, args...); err != nil {
		return nil, fmt.Errorf("select user: %w", err)
	}

	return &user, nil
}

func (ur *userRepository) ChangeName(ctx context.Context, newName string, userID uuid.UUID) error {
	sql, args, err := ur.qb.Builder.Update("users").Set("name", newName).
		Where(squirrel.Eq{"id": userID}).ToSql()
	if err != nil {
		return ErrFailBuildQuery
	}

	result, err := ur.db.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("update name: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return ErrGetAffected
	}

	if affected == 0 {
		return ErrInvalidUserID
	}

	return nil
}

func (ur *userRepository) ChangeEmail(ctx context.Context, newEmail string, userID uuid.UUID) error {
	sql, args, err := ur.qb.Builder.Update("users").Set("email", newEmail).
		Where(squirrel.Eq{"id": userID}).ToSql()
	if err != nil {
		return ErrFailBuildQuery
	}

	result, err := ur.db.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("update email: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return ErrGetAffected
	}

	if affected == 0 {
		return ErrInvalidUserID
	}

	return nil
}

func (ur *userRepository) ChangePassword(ctx context.Context, newPassword string, userID uuid.UUID) error {
	sql, args, err := ur.qb.Builder.Update("users").Set("password", newPassword).
		Where(squirrel.Eq{"id": userID}).ToSql()
	if err != nil {
		return ErrFailBuildQuery
	}

	result, err := ur.db.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return ErrGetAffected
	}

	if affected == 0 {
		return ErrInvalidUserID
	}

	return nil
}

func (ur *userRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	sql, args, err := ur.qb.Builder.Delete("users").
		Where(squirrel.Eq{"id": userID}).ToSql()
	if err != nil {
		return ErrFailBuildQuery
	}

	result, err := ur.db.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return ErrGetAffected
	}

	if affected == 0 {
		return ErrInvalidUserID
	}

	return nil
}
