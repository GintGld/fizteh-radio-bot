package service

import "errors"

var (
	// Auth
	ErrUserNotFound = errors.New("user not found")

	// Links
	ErrInvalidLink = errors.New("invalid link")
)
