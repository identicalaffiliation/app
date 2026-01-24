package rest

import "errors"

var (
	ErrInvalidMethod   error = errors.New("invalid http method")
	ErrInvalidJSONBody error = errors.New("invalid JSON body")
)
