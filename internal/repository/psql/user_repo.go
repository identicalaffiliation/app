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
		ur.logger.Logger.Error("failed to build query for create user",
			"operation", "create user",
			"user_id", user.ID.String(),
			"error", err.Error(),
		)

		return ErrFailBuildQuery
	}
	err = ur.db.DB.QueryRowxContext(ctx, sql, args...).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		ur.logger.Logger.Error("failed to create user",
			"operation", "create user",
			"user_id", user.ID.String(),
			"error", err.Error(),
		)

		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

func (ur *userRepository) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	sql, args, err := ur.qb.Builder.Select("id", "name", "email", "password",
		"created_at", "updated_at").From("users").OrderBy("email").ToSql()
	if err != nil {
		ur.logger.Logger.Error("failed to build query for get users",
			"operation", "get users",
			"error", err.Error(),
		)

		return nil, ErrFailBuildQuery
	}

	var users []*entity.User
	if err := ur.db.DB.SelectContext(ctx, &users, sql, args...); err != nil {
		ur.logger.Logger.Error("failed to get users",
			"operation", "get users",
			"error", err.Error(),
		)

		return nil, fmt.Errorf("select users: %w", err)
	}

	return users, nil
}

func (ur *userRepository) GetByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	sql, args, err := ur.qb.Builder.Select("id, name, email, password, created_at, updated_at").
		From("users").Where(squirrel.Eq{"id": userID}).ToSql()
	if err != nil {
		ur.logger.Logger.Error("failed to build query for get user",
			"operation", "get user",
			"user_id", userID.String(),
			"error", err.Error(),
		)

		return nil, ErrFailBuildQuery
	}

	var user entity.User
	if err := ur.db.DB.GetContext(ctx, &user, sql, args...); err != nil {
		ur.logger.Logger.Error("failed to get user",
			"operation", "get user",
			"user_id", userID.String(),
			"error", err.Error(),
		)

		return nil, fmt.Errorf("select user: %w", err)
	}

	return &user, nil
}

func (ur *userRepository) ChangeName(ctx context.Context, newName string, userID uuid.UUID) error {
	sql, args, err := ur.qb.Builder.Update("users").Set("name", newName).
		Where(squirrel.Eq{"id": userID}).ToSql()
	if err != nil {
		ur.logger.Logger.Error("failed to build query for update name",
			"operation", "update name",
			"user_id", userID.String(),
			"req_name", newName,
			"error", err.Error(),
		)

		return ErrFailBuildQuery
	}

	result, err := ur.db.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		ur.logger.Logger.Error("failed to update name",
			"operation", "update name",
			"user_id", userID.String(),
			"req_name", newName,
			"error", err.Error(),
		)

		return fmt.Errorf("update name: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		ur.logger.Logger.Error("failed to get affected from update name",
			"operation", "update name",
			"user_id", userID.String(),
			"req_name", newName,
			"error", err.Error(),
		)

		return ErrGetAffected
	}

	if affected == 0 {
		ur.logger.Logger.Error("failed to update name",
			"operation", "update name",
			"user_id", userID.String(),
			"req_name", newName,
			"error", errors.New("user not found").Error(),
		)

		return errors.New("user not found")
	}

	return nil
}

func (ur *userRepository) ChangeEmail(ctx context.Context, newEmail string, userID uuid.UUID) error {
	sql, args, err := ur.qb.Builder.Update("users").Set("email", newEmail).
		Where(squirrel.Eq{"id": userID}).ToSql()
	if err != nil {
		ur.logger.Logger.Error("failed to build query for update email",
			"operation", "update email",
			"user_id", userID.String(),
			"req_email", newEmail,
			"error", err.Error(),
		)

		return ErrFailBuildQuery
	}

	result, err := ur.db.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		ur.logger.Logger.Error("failed to update email",
			"operation", "update email",
			"user_id", userID.String(),
			"req_email", newEmail,
			"error", err.Error(),
		)

		return fmt.Errorf("update email: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		ur.logger.Logger.Error("failed to get affected from update email",
			"operation", "update email",
			"user_id", userID.String(),
			"req_email", newEmail,
			"error", err.Error(),
		)

		return ErrGetAffected
	}

	if affected == 0 {
		ur.logger.Logger.Error("failed to update email",
			"operation", "update email",
			"user_id", userID.String(),
			"req_email", newEmail,
			"error", errors.New("user not found").Error(),
		)

		return errors.New("user not found")
	}

	return nil
}

func (ur *userRepository) ChangePassword(ctx context.Context, newPassword string, userID uuid.UUID) error {
	sql, args, err := ur.qb.Builder.Update("users").Set("password", newPassword).
		Where(squirrel.Eq{"id": userID}).ToSql()
	if err != nil {
		ur.logger.Logger.Error("failed to build query for update password",
			"operation", "update password",
			"user_id", userID.String(),
			"req_password", newPassword,
			"error", err.Error(),
		)

		return ErrFailBuildQuery
	}

	result, err := ur.db.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		ur.logger.Logger.Error("failed to update password",
			"operation", "update password",
			"user_id", userID.String(),
			"req_password", newPassword,
			"error", err.Error(),
		)

		return fmt.Errorf("update password: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		ur.logger.Logger.Error("failed to get affected from update password",
			"operation", "update password",
			"user_id", userID.String(),
			"req_password", newPassword,
			"error", err.Error(),
		)

		return ErrGetAffected
	}

	if affected == 0 {
		ur.logger.Logger.Error("failed to update password",
			"operation", "update password",
			"user_id", userID.String(),
			"req_password", newPassword,
			"error", errors.New("user not found").Error(),
		)

		return errors.New("user not found")
	}

	return nil
}

func (ur *userRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	sql, args, err := ur.qb.Builder.Delete("users").
		Where(squirrel.Eq{"id": userID}).ToSql()
	if err != nil {
		ur.logger.Logger.Error("failed to build query for delete user",
			"operation", "delete user",
			"user_id", userID.String(),
			"error", err.Error(),
		)

		return ErrFailBuildQuery
	}

	result, err := ur.db.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		ur.logger.Logger.Error("failed to delete user",
			"operation", "delete user",
			"user_id", userID.String(),
			"error", err.Error(),
		)

		return fmt.Errorf("delete user: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		ur.logger.Logger.Error("failed to get affected from delete user",
			"operation", "delete user",
			"user_id", userID.String(),
			"error", err.Error(),
		)

		return ErrGetAffected
	}

	if affected == 0 {
		ur.logger.Logger.Error("failed to delete user",
			"operation", "delete user",
			"user_id", userID.String(),
			"error", errors.New("user not found").Error(),
		)

		return errors.New("user not found")
	}

	return nil
}
