package psql

import (
	"context"

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
