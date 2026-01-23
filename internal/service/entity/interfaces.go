package entity

import (
	"context"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/app/internal/dto"
)

type AuthUseCases interface {
	Register(ctx context.Context, userRequest *dto.UserRegisterRequest) error
	Login(ctx context.Context, userRequest *dto.UserLoginRequest) (*dto.AuthResponse, error)
}

type UserUseCases interface {
	GetUser(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
	GetUsers(ctx context.Context) ([]*dto.UserResponse, error)
	ChangeName(ctx context.Context, changeNameRequest *dto.ChangeUserNameRequest) error
	ChangeEmail(ctx context.Context, changeEmailRequest *dto.ChangeUserEmailRequest) error
	ChangePassword(ctx context.Context, changePasswordRequest *dto.ChangeUserPasswordRequest) error
	// DeleteUser(ctx context.Context, userID uuid.UUID) error
}
