package client

import "errors"

var (
	ErrTrackNotFound = errors.New("track not found")

	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrNotAuthorized       = errors.New("not authorized")
	ErrInternalServerError = errors.New("internal server error")
)
