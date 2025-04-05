package errors

import "errors"

var (
	ErrTimeout = errors.New("request timeout")
	ErrNetwork = errors.New("network unreachable")
)
