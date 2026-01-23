package entity

import "errors"

var (
	ErrInvalidUserEmail error = errors.New("invalid user email")
	ErrInvalidPassword  error = errors.New("invalid password")
	ErrInvalidUserID    error = errors.New("invalid user ID")
)
