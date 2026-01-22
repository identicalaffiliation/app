package dto

import "github.com/google/uuid"

type UserCreateRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type ChangeUserNameRequest struct {
	ID       uuid.UUID `validate:"required"`
	Name     string    `json:"name" validate:"required,min=2"`
	Password string    `json:"password" validate:"required,min=8"`
}

type ChangeUserEmailRequest struct {
	ID       uuid.UUID `validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8"`
}

type ChangeUserPasswordRequest struct {
	ID          uuid.UUID `validate:"required"`
	OldPassword string    `json:"oldPassword" validate:"required,min=8"`
	NewPassword string    `json:"NewPassword" validate:"required,min=8"`
}
