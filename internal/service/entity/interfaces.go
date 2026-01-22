package entity

import (
	"context"

	"github.com/identicalaffiliation/app/internal/dto"
)

type UserUseCases interface {
	CreateUser(ctx context.Context, userRequest *dto.UserCreateRequest) error
	// ChangeName(ctx context.Context, changeRequest *dto.ChangeUserNameRequest) error
	// ChangeEmail(ctx context.Context, changeRequest *dto.ChangeUserEmailRequest) error
	// ChangePassword(ctx context.Context, changeRequest *dto.ChangeUserPasswordRequest) error
}
