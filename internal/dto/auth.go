package dto

import "time"

type UserRegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type AuthResponse struct {
	User      *UserResponse `json:"user"`
	Token     string        `json:"token"`
	ExpiresAt time.Time     `json:"expiresAt"`
}
