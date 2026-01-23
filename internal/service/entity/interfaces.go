package entity

import (
	"context"

	"github.com/identicalaffiliation/app/internal/dto"
)

type AuthUseCases interface {
	Register(ctx context.Context, userRequest *dto.UserRegisterRequest) error
	Login(ctx context.Context, userRequest *dto.UserLoginRequest) (*dto.AuthResponse, error)
}
