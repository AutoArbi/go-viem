package errors

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrInternalError  = errors.New("internal error")
)
