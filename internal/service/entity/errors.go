package entity

import "errors"

var (
	ErrInvalidCreateUserRequest error = errors.New("invalid create user request: %w")
)
