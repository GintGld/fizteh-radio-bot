package service

import "errors"

var (
	// Auth
	ErrUserNotFound = errors.New("user not found")

	ErrMediaExists = errors.New("media exists")

	ErrTagNotFound = errors.New("tag not found")

	// Links
	ErrInvalidLink = errors.New("invalid link")
)
