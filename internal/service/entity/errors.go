package entity

import "errors"

var (
	ErrInvalidUserEmail error = errors.New("invalid user email")
	ErrInvalidPassword  error = errors.New("invalid password")
	ErrInvalidUserID    error = errors.New("invalid user ID")

	ErrInvalidTodoStatus error = errors.New("cannot create todo that already have done")
	ErrInvalidTodoID     error = errors.New("invalid todo ID")
)
