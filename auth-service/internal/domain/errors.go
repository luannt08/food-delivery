package domain

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrProfileNotFound    = errors.New("profile not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)
