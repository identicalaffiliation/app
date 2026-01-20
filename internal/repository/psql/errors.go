package psql

import (
	"errors"
)

var (
	ErrFailBuildQuery error = errors.New("fail to build query")
	ErrInvalidUserID  error = errors.New("invalid user ID")
)
