package entity

import "errors"

var (
	ErrInvalidCreateUserRequest error = errors.New("invalid create user request: %w")
	ErrInvalidUserEmail         error = errors.New("invalid user email")
	ErrInvalidPassword          error = errors.New("invalid password")
)
