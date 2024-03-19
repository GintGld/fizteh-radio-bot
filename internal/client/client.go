package client

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrNotAuthorized       = errors.New("not authorized")
	ErrInternalServerError = errors.New("internal server error")
)
