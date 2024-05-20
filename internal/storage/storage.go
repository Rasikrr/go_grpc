package storage

import "errors"

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user bot found")
	ErrAppNotFound  = errors.New("app not found")
)
